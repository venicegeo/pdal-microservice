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

/*
pdal-microservice provides an endpoint for accepting PDAL requests.

Examples

  $ curl -v --noproxy hostIP -X POST -H "Content-Type: application/json" \
    -d '{"bucket":"venicegeo-sample-data","key":"pointcloud/samp11-utm.laz"}'
    http://hostIP:8080/info
*/
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

	// POST /info expects JSON specifying S3 bucket/key of the point cloud.
	// Alternatives should include a GRiD export primary key or URL. What about a local path?
	router.POST("/info", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		// Check that we have a valid path. Is this the correct place to do this?
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
			// errors here should also be JSON-encoded as below
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
			// errors here should also be JSON-encoded as below
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

		// output needs to be more meaningful
		t := msg //`{"status":"success"}`

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(t); err != nil {
			log.Fatal(err)
		}
	})

	fmt.Println("Starting up on 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		// errors here should also be JSON-encoded as above
		log.Fatal(err)
	}
}
