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

// Package servicecontroller provides pz-servicecontroller helper functions.
package servicecontroller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
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

// RegisterServiceMsg defines the expected output JSON returned by Piazza when
// an external service is registered.
type RegisterServiceMsg struct {
	ResourceID string `json:"resourceId"`
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
		return errors.New("Error marshaling ResourceMetadata")
	}

	response, err := http.Post(
		RegistryURL, ContentTypeJSON, bytes.NewBuffer(data),
	)
	if err != nil {
		return errors.New("Error posting ResourceMetadata to registerService")
	}

	if response.Body == nil {
		return errors.New("No JSON body returned from registerService")
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("Error reading JSON body returned from registerService")
	}

	// Throw 400 if we cannot unmarshal the body as a valid InputMsg.
	var rm RegisterServiceMsg
	if err := json.Unmarshal(b, &rm); err != nil {
		return errors.New("Error unmarshaling RegisterServiceMsg")
	}
	log.Println("RegisterService received resourceId=" + rm.ResourceID)

	return nil
}
