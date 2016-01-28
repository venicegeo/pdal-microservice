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

/*
pzsvc-pdal provides an endpoint for accepting PDAL requests.

Examples

  $ curl -v -X POST -H "Content-Type: application/json" \
    -d '{"source":{"bucket":"venicegeo-sample-data","key":"pointcloud/samp11-utm.laz"},"function":"info"}' http://hostIP:8080/pdal

We shall see where we land with the input and output message for the job manager, but for now, we are expecting something along these lines.

Input:

	{
		"source": {
			"bucket": "venicegeo-sample-data",
			"key": "pointcloud/samp11-utm.laz"
		},
		"function": "ground",
		"options": {
			"slope": 0.5
		},
		"destination": {
			"bucket": "venicegeo-sample-data",
			"key" "temp/output.laz"
		}
	}

Output:

	{
		"input": <echo the input message>,
		"started_at": "2015-12-23T18:07:36.987565884Z",
		"finished_at": "2015-12-23T18:07:38.111658707Z",
		"code": 200,
		"message": "Success!"
	}

These messages are known to be incomplete at the moment. I'm sure there will be things like job IDs, etc. that have not been included at the moment. This is a good starting point though.
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-pdal/functions"
	"github.com/venicegeo/pzsvc-pdal/handlers"
)

// var pdalMetadata = job.ResourceMetadata{
// 	Name:         "pdal",
// 	ServiceID:    "",
// 	Description:  "Process point cloud data using PDAL.",
// 	URL:          "https://api.piazzageo.io/v1/pdal",
// 	Networks:     "TBD",
// 	QoS:          "Development",
// 	Availability: "UP",
// 	Tags:         "point cloud, pdal, lidar",
// 	ClassType:    "Unclassified",
// 	// TermDate:            time.Now(),
// 	// ClientCertRequired:  false,
// 	// CredentialsRequired: false,
// 	// PreAuthRequired:     false,
// 	// Contracts:           "",
// 	// Method:              "",
// 	// MimeType:            "",
// 	// Params:              "",
// 	// Reason:              "",
// }

func main() {
	// For standalone demo purposes, we will start two services: our PDAL service, and a mocked up JobManager.

	router := httprouter.New()

	router.GET("/",
		func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			fmt.Fprintf(w, "Hi!")
		})

	type ListFuncs struct {
		Functions []string `json:"functions"`
	}
	out := ListFuncs{[]string{
		"crop", "dart", "dtm", "ground", "height", "info", "radius", "statistical",
		"translate", "vo",
	}}

	router.GET("/functions/:name",
		func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			var a interface{}
			switch ps.ByName("name") {
			case "crop":
				a = functions.NewCropOptions()
				w.WriteHeader(http.StatusOK)

			case "dart":
				a = functions.NewDartOptions()
				w.WriteHeader(http.StatusOK)

			case "dtm":
				a = functions.NewDtmOptions()
				w.WriteHeader(http.StatusOK)

			case "ground":
				a = functions.NewGroundOptions()
				w.WriteHeader(http.StatusOK)

			case "height":
				w.WriteHeader(http.StatusOK)

			case "info":
				a = functions.NewInfoOptions()
				w.WriteHeader(http.StatusOK)

			case "radius":
				a = functions.NewRadiusOptions()
				w.WriteHeader(http.StatusOK)

			case "statistical":
				a = functions.NewStatisticalOptions()
				w.WriteHeader(http.StatusOK)

			case "translate":
				w.WriteHeader(http.StatusOK)

			case "vo":
				w.WriteHeader(http.StatusOK)

			default:
				type DefaultMsg struct {
					Message string `json:"message"`
					ListFuncs
				}
				msg := "Unrecognized function " + ps.ByName("name") + "."
				a = DefaultMsg{msg, out}
				w.WriteHeader(http.StatusBadRequest)
			}
			if err := json.NewEncoder(w).Encode(a); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

	router.GET("/functions",
		func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(out); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

	// // Setup the PDAL service.
	router.POST("/pdal", handlers.PdalHandler)

	// Setup the mocked up JobManager.
	// router.POST("/manager", handlers.JobManagerHandler)

	var defaultPort = os.Getenv("PORT")
	if defaultPort == "" {
		defaultPort = "8080"
	}

	log.Println("Starting on port ", defaultPort)
	log.Println(os.Getenv("PATH"))
	log.Println(os.Getenv("LD_LIBRARY_PATH"))
	if err := http.ListenAndServe(":"+defaultPort, router); err != nil {
		log.Fatal(err)
	}
}
