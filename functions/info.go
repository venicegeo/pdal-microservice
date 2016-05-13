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
	"errors"
	"log"
	"os/exec"
)

// InfoOptions defines options for the Info function.
type InfoOptions struct {
	// compute hexagonal boundary that contains all points
	Boundary bool `json:"boundary"`
	// dump metadata associated with the input file
	Metadata bool `json:"metadata"`
	// dump the schema of the internal point storage
	Schema bool `json:"schema"`
}

// NewInfoOptions constructs InfoOptions with default values.
func NewInfoOptions() *InfoOptions {
	return &InfoOptions{
		Boundary: false,
		Metadata: false,
		Schema:   false,
	}
}

// Info implements pdal info.
func Info(i, o string, options *json.RawMessage) ([]byte, error) {
	opts := NewInfoOptions()
	if options != nil {
		if err := json.Unmarshal(*options, &opts); err != nil {
			return nil, errors.New("Error with json.Unmarshal() " + err.Error())
		}
	}
	log.Println("detected info options")

	var args []string
	args = append(args, "info")
	args = append(args, i)
	if opts.Boundary {
		args = append(args, "--boundary")
	}
	if opts.Metadata {
		args = append(args, "--metadata")
	}
	if opts.Schema {
		args = append(args, "--schema")
	}
	log.Println("set args " + args...)

	out, err := exec.Command("pdal", args...).CombinedOutput()
	if err != nil {
		return nil, errors.New("Error with exec.Command() " + err.Error())
	}
	log.Println("ran the command")

	// Trim whitespace
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, out); err != nil {
		return nil, errors.New("Error with json.Compact() " + err.Error())
	}
	log.Println("compact buffer")

	return buffer.Bytes(), nil
}
