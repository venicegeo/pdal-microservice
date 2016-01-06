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

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/venicegeo/pzsvc-sdk-go/objects"
)

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

/*
BadRequest handles bad requests.

All bad requests result in a failure in the eyes of the JobManager. The ResponseWriter echos some key aspects of the Request (e.g., input, start time) and appends StatusBadRequest (400) as well as a message to the JobOutput, which is returned as JSON.
*/
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

/*
InternalError handles internal server errors.

All internal server errors result in an error in the eyes of the JobManager. The ResponseWriter echos some key aspects of the Request (e.g., input, start time) and appends StatusInternalServerError (500) as well as a message to the JobOutput, which is returned as JSON.
*/
func InternalError(w http.ResponseWriter, r *http.Request, res *objects.JobOutput, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	res.Code = http.StatusInternalServerError
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
	UpdateJobManager(objects.Error, r)
}

/*
Okay handles successful calls.

All successful calls result in sucess in the eyes of the JobManager. The ResponseWriter echos some key aspects of the Request (e.g., input, start time) and appends StatusOK (200) as well as a message to the JobOutput, which is returned as JSON.
*/
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

/*
S3Download downloads a file from an S3 bucket/key.
*/
func S3Download(file *os.File, bucket, key string) error {
	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	log.Println("Downloaded", numBytes, "bytes")
	return nil
}

/*
S3Upload uploads a file to an S3 bucket.
*/
func S3Upload(file *os.File, bucket, key string) error {
	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	log.Println("Successfully uploaded to", result.Location)
	return nil
}

// GetJobInput provides a common means of parsing the JobInput JSON.
func GetJobInput(w http.ResponseWriter, r *http.Request, res objects.JobOutput) objects.JobInput {
	var msg objects.JobInput

	// There should always be a body, else how are we to know what to do? Throw
	// 400 if missing.
	if r.Body == nil {
		BadRequest(w, r, res, "No JSON")
	}

	// Throw 500 if we cannot read the body.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		InternalError(w, r, &res, err.Error())
	}

	// Throw 400 if we cannot unmarshal the body as a valid JobInput.
	if err := json.Unmarshal(b, &msg); err != nil {
		BadRequest(w, r, res, err.Error())
	}

	return msg
}

// ParseFilenameFromKey parses the S3 filename from the key.
func ParseFilenameFromKey(key string) string {
	keySlice := strings.Split(key, "/")
	return keySlice[len(keySlice)-1]
}

// FunctionFunc defines the signature of our function creator.
type FunctionFunc func(http.ResponseWriter, *http.Request,
	*objects.JobOutput, objects.JobInput)

// MakeFunction wraps the individual PDAL functions.
// Parse the input and output filenames, creating files as needed. Download the
// input data and upload the output data.
func MakeFunction(fn func(http.ResponseWriter, *http.Request,
	*objects.JobOutput, objects.JobInput, string, string)) FunctionFunc {
	return func(w http.ResponseWriter, r *http.Request, res *objects.JobOutput,
		msg objects.JobInput) {
		var inputName, outputName string
		var fileIn, fileOut *os.File

		// Split the source S3 key string, interpreting the last element as the
		// input filename. Create the input file, throwing 500 on error.
		inputName = ParseFilenameFromKey(msg.Source.Key)
		fileIn, err := os.Create(inputName)
		if err != nil {
			InternalError(w, r, res, err.Error())
			return
		}
		defer fileIn.Close()

		// If provided, split the destination S3 key string, interpreting the last
		// element as the output filename. Create the output file, throwing 500 on
		// error.
		if len(msg.Destination.Key) > 0 {
			outputName = ParseFilenameFromKey(msg.Destination.Key)
			fileOut, err = os.Create(outputName)
			if err != nil {
				InternalError(w, r, res, err.Error())
				return
			}
			defer fileOut.Close()
		}

		// Download the source data from S3, throwing 500 on error.
		err = S3Download(fileIn, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			InternalError(w, r, res, err.Error())
			return
		}

		// Run the PDAL function.
		fn(w, r, res, msg, inputName, outputName)

		// If an output has been created, upload the destination data to S3,
		// throwing 500 on error.
		if len(msg.Destination.Key) > 0 {
			err = S3Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
			if err != nil {
				InternalError(w, r, res, err.Error())
				return
			}
		}
	}
}
