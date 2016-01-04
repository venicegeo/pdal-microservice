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

package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/venicegeo/pzsvc-pdal/objects"
	"github.com/venicegeo/pzsvc-pdal/utils"
)

// infoOptions defines options for Info.
type infoOptions struct {
	Boundary *bool `json:"boundary"`
	Metadata *bool `json:"metadata"`
	Schema   *bool `json:"schema"`
}

// InfoFunction implements pdal info.
func InfoFunction(w http.ResponseWriter, r *http.Request,
	res *objects.JobOutput, msg objects.JobInput, i, o string) {
	boundary := false
	metadata := false
	schema := false
	if msg.Options != nil {
		var opts infoOptions
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			utils.BadRequest(w, r, *res, err.Error())
			return
		}
		if opts.Boundary != nil {
			boundary = *opts.Boundary
		}
		if opts.Metadata != nil {
			metadata = *opts.Metadata
		}
		if opts.Schema != nil {
			schema = *opts.Schema
		}
	}
	var params string
	if boundary {
		params = params + "--boundary"
	}
	if metadata {
		params = params + "--metadata"
	}
	if schema {
		params = params + "--schema"
	}

	out, _ := exec.Command("pdal", *msg.Function, i, params).CombinedOutput()

	// Trim whitespace
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, out); err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(buffer.Bytes(), &res.Response); err != nil {
		log.Fatal(err)
	}
}
