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

// Package job provides JobManager helper functions.
package job

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// StatusType is a string describing the state of the job.
type StatusType int

// Enumerate valid StatusType values.
const (
	Submitted StatusType = iota
	Running
	Success
	Cancelled
	Error
	Fail
)

var statuses = [...]string{
	"submitted", "running", "success", "cancelled", "error", "fail",
}

func (status StatusType) String() string {
	return statuses[status]
}

// S3Bucket defines the expected JSON structure for S3 buckets.
// An S3 bucket can be used for source (input) and destination (output) files.
type S3Bucket struct {
	Bucket string `json:"bucket,omitempty"`
	Key    string `json:"key,omitempty"`
}

// InputMsg defines the expected input JSON structure.
// We currently support S3 input (bucket/key), though provider-specific (e.g.,
// GRiD) may be legitimate.
type InputMsg struct {
	Source      S3Bucket         `json:"source,omitempty"`
	Function    *string          `json:"function,omitempty"`
	Options     *json.RawMessage `json:"options,omitempty"`
	Destination S3Bucket         `json:"destination,omitempty"`
}

// OutputMsg defines the expected output JSON structure.
type OutputMsg struct {
	Input      InputMsg                    `json:"input,omitempty"`
	StartedAt  time.Time                   `json:"started_at,omitempty"`
	FinishedAt time.Time                   `json:"finished_at,omitempty"`
	Code       int                         `json:"code,omitempty"`
	Message    string                      `json:"message,omitempty"`
	Response   map[string]*json.RawMessage `json:"response,omitempty"`
}

// ResourceMetadata defines the metadata required to register the service.
type ResourceMetadata struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	URL              string `json:"url"`
	Method           string `json:"method,omitempty"`
	RequestMimeType  string `json:"requestMimeType,omitempty"`
	ResponseMimeType string `json:"responseMimeType,omitempty"`
	Params           string `json:"params,omitempty"`
}

// RegisterServiceMsg defines the expected output JSON returned by Piazza when
// an external service is registered.
type RegisterServiceMsg struct {
	ResourceID string `json:"resourceId"`
}

// UpdateMsg defines the expected output JSON structure for updating the
// JobManager.
type UpdateMsg struct {
	Status string `json:"status"`
}

// Update handles PDAL status updates.
func Update(t StatusType, r *http.Request) {
	log.Println("Setting job status as \"", t.String(), "\"")
	// var res UpdateMsg
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

All bad requests result in a failure in the eyes of the JobManager. The
ResponseWriter echos some key aspects of the Request (e.g., input, start time)
and appends StatusBadRequest (400) as well as a message to the OutputMsg, which
is returned as JSON.
*/
func BadRequest(w http.ResponseWriter, r *http.Request, res OutputMsg, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	res.Code = http.StatusBadRequest
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	Update(Fail, r)
}

/*
InternalError handles internal server errors.

All internal server errors result in an error in the eyes of the JobManager. The
ResponseWriter echos some key aspects of the Request (e.g., input, start time)
and appends StatusInternalServerError (500) as well as a message to the
OutputMsg, which is returned as JSON.
*/
func InternalError(
	w http.ResponseWriter, r *http.Request, res OutputMsg, message string,
) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	res.Code = http.StatusInternalServerError
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	Update(Error, r)
}

/*
Okay handles successful calls.

All successful calls result in sucess in the eyes of the JobManager. The
ResponseWriter echos some key aspects of the Request (e.g., input, start time)
and appends StatusOK (200) as well as a message to the OutputMsg, which is
returned as JSON.
*/
func Okay(
	w http.ResponseWriter, r *http.Request, res OutputMsg, message string,
) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	res.Code = http.StatusOK
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	Update(Success, r)
}

// GetInputMsg provides a common means of parsing the InputMsg JSON.
func GetInputMsg(
	w http.ResponseWriter, r *http.Request, res OutputMsg,
) InputMsg {
	var msg InputMsg

	// There should always be a body, else how are we to know what to do? Throw
	// 400 if missing.
	if r.Body == nil {
		http.Error(w, "No JSON", http.StatusBadRequest)
		return msg
	}

	// Throw 500 if we cannot read the body.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return msg
	}

	// Throw 400 if we cannot unmarshal the body as a valid InputMsg.
	if err := json.Unmarshal(b, &msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return msg
	}

	return msg
}

// ContentTypeJSON is the http content-type for JSON.
const ContentTypeJSON = "application/json"

// registryURL is the Piazza registration endpoint
const RegistryURL = "http://pz-servicecontroller.cf.piazzageo.io/servicecontroller/registerService"

//const RegistryURL = "http://localhost:8082/servicecontroller/registerService"

/*
RegisterService handles service registartion with Piazza for external services.
*/
func RegisterService(m ResourceMetadata) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	response, err := http.Post(RegistryURL, ContentTypeJSON, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if response.Body == nil {
		return errors.New("No JSON body returned from registerService")
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Throw 400 if we cannot unmarshal the body as a valid InputMsg.
	var rm RegisterServiceMsg
	if err := json.Unmarshal(b, &rm); err != nil {
		return err
	}
	log.Println("RegisterService received resourceId=" + rm.ResourceID)

	return nil
}
