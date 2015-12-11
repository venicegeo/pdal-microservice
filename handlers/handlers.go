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

// UpdateDispatcher handles PDAL status updates.
func UpdateDispatcher(w http.ResponseWriter, t objects.StatusType) {
	log.Println("Setting job status as a", t.String())
	var res objects.DispatcherUpdate
	res.Status = t.String()
	// update the dispatcher/job table
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
}

// PdalHandler handles PDAL jobs.
func PdalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Received request")
	var res objects.JobOutput
	res.StartedAt = time.Now()

	// Check that we have a valid path. Is this the correct place to do this?
	log.Println("Checking to see if", r.URL.Path, "is a valid endpoint")
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	log.Println("Attempt to read the JSON body")
	// Parse the incoming JSON body, and unmarshal as events.NewData struct.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		UpdateDispatcher(w, objects.Error)
		log.Fatal(err)
	}

	log.Println("Attempt to unmarshal the JSON")
	var msg objects.JobInput
	if err := json.Unmarshal(b, &msg); err != nil {
		UpdateDispatcher(w, objects.Fail)
		log.Fatal(err)
	}
	if msg.Function == nil {
		UpdateDispatcher(w, objects.Fail)
		log.Println("Must provide a function")
		return
	}

	res.Input = msg
	res.Status = objects.Running.String()
	// we have successfully parsed the input JSON, update dispatcher/job table that we are now running

	file, err := os.Create("download_file.laz")
	if err != nil {
		// errors here should also be JSON-encoded as below
		http.Error(w, err.Error(), http.StatusInternalServerError)
		res.Status = objects.Error.String()
		// update the dispatcher/job table
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
		// errors here should also be JSON-encoded as below
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
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
	res.Status = objects.Success.String()
	res.FinishedAt = time.Now()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
}