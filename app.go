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
    -d '{"source":{"bucket":"venicegeo-sample-data","key":"pointcloud/samp11-utm.laz"},"function":"info"}' http://hostIP:8080/pdal
*/
package main

import (
	"log"
	"net/http"

	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pdal-microservice/handlers"
)

func main() {
	router := httprouter.New()

	router.POST("/pdal", handlers.PdalHandler)

	log.Println("Starting /pdal on 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		// errors here should also be JSON-encoded as above
		log.Fatal(err)
	}
}
