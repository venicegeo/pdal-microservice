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

// Package gateway provides pz-gateway helper functions.
package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

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

// JobMsg defines the expected output JSON returned by Piazza when
// an external service is registered.
type JobMsg struct {
	Type  string `json:"type"`
	JobID string `json:"jobId"`
}

// GatewayJobURL is the Piazza registration endpoint
const GatewayJobURL = "http://pz-gateway.cf.piazzageo.io/job"

/*
RegisterService handles service registartion with Piazza for external services.
*/
func RegisterService(m ResourceMetadata) error {
	data, err := json.Marshal(m)
	if err != nil {
		return errors.New("Error marshaling ResourceMetadata")
	}

	str := fmt.Sprintf("{\"apiKey\":\"my-api-key-38n987\",\"jobType\":{\"type\":\"register-service\",\"data\":%s}}", bytes.NewBuffer(data))
	fmt.Println(str)

	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)
	err = w.WriteField("body", str)
	if err != nil {
		return errors.New("Error writing body")
	}
	err = w.Close()
	if err != nil {
		return errors.New("Error closing writer")
	}

	req, err := http.NewRequest("POST", GatewayJobURL, &buffer)
	if err != nil {
		return errors.New("Error creating request")
	}
	req.Header.Add("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return errors.New("Error performing request")
	}

	if response.Body == nil {
		return errors.New("No JSON body returned from gateway")
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("Error reading JSON body returned from gateway")
	}

	var rm JobMsg
	if err := json.Unmarshal(b, &rm); err != nil {
		return errors.New("Error unmarshaling JobMsg")
	}
	log.Println("Gateway received type " + rm.Type + ", jobId " + rm.JobID)

	return nil
}
