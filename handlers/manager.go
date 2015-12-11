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
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pdal-microservice/objects"
)

// JobManagerHandler handles JobManager updates.
func JobManagerHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var msg objects.JobManagerUpdate
	if err := json.Unmarshal(b, &msg); err != nil {
		log.Fatal(err)
	}

	log.Println("Job status reported as \"", msg.Status, "\"")
}