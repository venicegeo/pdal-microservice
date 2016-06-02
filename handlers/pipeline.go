/*
Copyright 2015-2016, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var readerNum int

func inferAndDownload(rawurl *url.URL) (name string, size int64, err error) {
	switch rawurl.Scheme {
	case "s3":
		bucket := rawurl.Host
		key := strings.TrimLeft(rawurl.Path, "/")
		ext := filepath.Ext(key)

		fname := "download_file-" + strconv.Itoa(readerNum) + ext
		readerNum++
		file, err := os.Create(fname)
		if err != nil {
			return "", 0, err
		}
		defer file.Close()

		downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
		numBytes, err := downloader.Download(file, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return "", 0, err
		}
		return file.Name(), numBytes, nil
	case "http", "https":
		client := &http.Client{}

		req, err := http.NewRequest("GET", rawurl.String(), nil)

		resp, err := client.Do(req)
		if err != nil {
			return "", 0, err
		}
		defer resp.Body.Close()
		ext := filepath.Ext(rawurl.String())

		fname := "download_file-" + strconv.Itoa(readerNum) + ext
		readerNum++
		file, err := os.Create(fname)
		if err != nil {
			return "", 0, err
		}
		defer file.Close()

		numBytes, err := io.Copy(file, resp.Body)
		if err != nil {
			return "", 0, err
		}
		return file.Name(), numBytes, nil
	default:
		log.Println("other")
		_, err := os.Stat(rawurl.String())
		if err == nil {
			log.Println("File exists")
		} else if os.IsNotExist(err) {
			log.Println("File does not exist")
		} else {
			return "", 0, err
		}
	}
	return "", 0, nil
}

func inferAndUpload(file *os.File, rawurl *url.URL) (location string, err error) {
	switch rawurl.Scheme {
	case "s3":
		bucket := rawurl.Host
		key := strings.TrimLeft(rawurl.Path, "/")

		uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
		result, err := uploader.Upload(&s3manager.UploadInput{
			Body:   file,
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return "", err
		}
		return result.Location, nil
	case "http", "https":
		log.Println("Not handled yet")
	default:
		log.Println("other")
		_, err := os.Stat(rawurl.String())
		if err == nil {
			log.Println("File exists")
		} else if os.IsNotExist(err) {
			log.Println("File does not exist")
		} else {
			return "", err
		}
	}
	return "", nil
}

// PipelineHandler handles PDAL jobs.
func PipelineHandler(w http.ResponseWriter, r *http.Request) *AppError {
	readerNum = 0

	// There should always be a body, else how are we to know what to do? Throw
	// 400 if missing.
	if r.Body == nil {
		return &AppError{nil, "No JSON", http.StatusBadRequest}
	}

	// Throw 500 if we cannot read the body.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	var uploadurl *url.URL

	var compactInput bytes.Buffer
	err = json.Compact(&compactInput, b)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}
	log.Println("Received input pipeline of:", compactInput.String())

	type pipeline struct {
		Pipeline []json.RawMessage `json:"pipeline"`
	}

	var p pipeline
	err = json.Unmarshal(b, &p)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}
	if p.Pipeline == nil {
		return &AppError{err, "No pipeline key", http.StatusInternalServerError}
	}

	var q pipeline

	numStages := len(p.Pipeline)

	for k, v := range p.Pipeline {
		var compactStage bytes.Buffer
		err := json.Compact(&compactStage, v)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}
		log.Println("Stage", k+1, "of", numStages, "parsed as:", compactStage.String())

		var vv string
		err = json.Unmarshal(v, &vv)
		if err != nil {
			q.Pipeline = append(q.Pipeline, v)
		} else {
			if k < numStages-1 {
				s3url, err := url.Parse(vv)
				if err != nil {
					return &AppError{err, err.Error(), http.StatusInternalServerError}
				}

				name, numBytes, err := inferAndDownload(s3url)
				if err != nil {
					return &AppError{err, err.Error(), http.StatusInternalServerError}
				}
				log.Println("Downloaded", name, numBytes, "bytes")
				q.Pipeline = append(q.Pipeline, json.RawMessage("\""+name+"\""))
			} else {
				uploadurl, err = url.Parse(vv)
				if err != nil {
					return &AppError{err, err.Error(), http.StatusInternalServerError}
				}

				q.Pipeline = append(q.Pipeline, json.RawMessage("\"upload_file.laz\""))
			}
		}
	}
	b, err = json.Marshal(q)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	var compactModified bytes.Buffer
	err = json.Compact(&compactModified, b)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}
	log.Println("Modified pipeline to:", compactModified.String())

	pipe, err := os.Create("pipeline.json")
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}
	defer pipe.Close()
	n, err := compactModified.WriteTo(pipe)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}
	log.Println("Wrote", n, "bytes to", pipe.Name())

	var args []string
	args = append(args, "pipeline")
	args = append(args, pipe.Name())
	args = append(args, "-v", "10", "--debug")

	outcmd, err := exec.Command("pdal", args...).CombinedOutput()

	log.Println(string(outcmd))
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	file, err := os.Open("upload_file.laz")
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}
	defer file.Close()

	location, err := inferAndUpload(file, uploadurl)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}
	log.Println("Uploaded to", location)

	type locationResult struct {
		Location string `json:"location"`
	}
	result := locationResult{Location: location}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	return nil
}
