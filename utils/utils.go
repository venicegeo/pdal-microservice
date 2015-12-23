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

package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/venicegeo/pdal-microservice/objects"
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
