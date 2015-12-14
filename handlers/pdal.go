/*
Copyright 2015, RadiantBlue Technologies, Inc.

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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pdal-microservice/objects"
)

// var validPath = regexp.MustCompile("^/(info|pipeline)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(pdal)$")

// UpdateJobManager handles PDAL status updates.
func UpdateJobManager(t objects.StatusType, r *http.Request) {
	log.Println("Setting job status as \"", t.String(), "\"")
	// var res objects.JobManagerUpdate
	// res.Status = t.String()
	// //	url := "http://192.168.99.100:8080/manager"
	// url := r.URL.Path + `/manager`
	//
	// jsonStr, err := json.Marshal(res)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("Content-Type", "application/json")
	//
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
}

// BadRequest handles bad requests.
func BadRequest(w http.ResponseWriter, r *http.Request, res objects.JobOutput, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	res.Code = http.StatusBadRequest
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
	UpdateJobManager(objects.Fail, r)
}

// InternalError handles internal server errors.
func InternalError(w http.ResponseWriter, r *http.Request, res objects.JobOutput, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	res.Code = http.StatusInternalServerError
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
	UpdateJobManager(objects.Error, r)
}

// Okay handles good requests.
func Okay(w http.ResponseWriter, r *http.Request, res objects.JobOutput, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	res.Code = http.StatusOK
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
	UpdateJobManager(objects.Success, r)
}

// PdalHandler handles PDAL jobs.
func PdalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var res objects.JobOutput
	res.StartedAt = time.Now()

	// Check that we have a valid path. Is this the correct place to do this?
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		BadRequest(w, r, res, "Endpoint does not exist")
		return
	}

	if r.Body == nil {
		BadRequest(w, r, res, "No JSON")
		return
	}

	// Parse the incoming JSON body, and unmarshal as events.NewData struct.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		InternalError(w, r, res, err.Error())
		return
	}

	var msg objects.JobInput
	if err := json.Unmarshal(b, &msg); err != nil {
		BadRequest(w, r, res, err.Error())
		return
	}
	if msg.Function == nil {
		BadRequest(w, r, res, "Must provide a function")
		return
	}

	log.Println("/pdal processing the data in", msg.Source.Bucket, "/", msg.Source.Key, "with", *msg.Function)

	res.Input = msg
	UpdateJobManager(objects.Running, r)

	file, err := os.Create("download_file.laz")
	if err != nil {
		InternalError(w, r, res, err.Error())
		return
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(msg.Source.Bucket),
			Key:    aws.String(msg.Source.Key),
		})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		InternalError(w, r, res, err.Error())
		return
	}
	log.Println("Downloaded", numBytes, "bytes")

	out, _ := exec.Command("pdal", *msg.Function, file.Name()).CombinedOutput()

	// Trim whitespace
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, out); err != nil {
		fmt.Println(err)
	}

	if err = json.Unmarshal(buffer.Bytes(), &res.Response); err != nil {
		log.Fatal(err)
	}

	res.FinishedAt = time.Now()
	Okay(w, r, res, "Success!")
}
