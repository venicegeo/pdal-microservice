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
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-pdal/handlers"
)

func main() {
	// For standalone demo purposes, we will start two services: our PDAL service, and a mocked up JobManager.

	router := httprouter.New()

	// Setup the PDAL service.
	router.POST("/pdal", handlers.PdalHandler)

	// Setup the mocked up JobManager.
	router.POST("/manager", handlers.JobManagerHandler)

	log.Println("Starting on 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
