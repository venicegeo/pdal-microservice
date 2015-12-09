// Copyright 2015, RadiantBlue Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
)

// var validPath = regexp.MustCompile("^/(info|pipeline)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(info)$")

func main() {
	router := httprouter.New()
	router.POST("/info", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		// Parse the incoming JSON body, and unmarshal as events.NewData struct.
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		type SourceBucket struct {
			Bucket string `json:"bucket"`
			Key    string `json:"key"`
		}
		var msg SourceBucket
		if err := json.Unmarshal(b, &msg); err != nil {
			log.Fatal(err)
		}

		file, err := os.Create("download_file.laz")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
		numBytes, err := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(msg.Bucket),
				Key:    aws.String(msg.Key),
			})
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				log.Println("Error:", awsErr.Code(), awsErr.Message())
			} else {
				fmt.Println(err.Error())
			}
			return
		}

		fmt.Fprintln(w, "Downloaded file", file.Name(), numBytes, "bytes")

		out, _ := exec.Command("pdal", "info", file.Name()).CombinedOutput()
		fmt.Fprintln(w, string(out))
	})

	fmt.Println("Starting up on 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
