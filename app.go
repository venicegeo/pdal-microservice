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
pzsvc-pdal provides an endpoint for accepting PDAL requests.

Examples

  $ curl -v -X POST -H "Content-Type: application/json" \
    -d '{"source":{"bucket":"venicegeo-sample-data","key":"pointcloud/samp11-utm.laz"},"function":"info"}' http://hostIP:8080/pdal
*/
package main

import (
	"log"
	"net/http"

	"github.com/venicegeo/pzsvc-pdal/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
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
