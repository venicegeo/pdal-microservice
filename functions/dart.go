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
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

// DartOptions defines options for the Dart function.
type DartOptions struct {
	Radius float64 `json:"radius"` // minimum distance between samples
}

// NewDartOptions constructs DartOptions with default values.
func NewDartOptions() *DartOptions {
	return &DartOptions{Radius: 1.0}
}

// Dart implements pdal height.
func Dart(i, o string, options *json.RawMessage) ([]byte, error) {
	opts := NewDartOptions()
	if options != nil {
		if err := json.Unmarshal(*options, &opts); err != nil {
			return nil, err
		}
	}

	var args []string
	args = append(args, "translate", i, o, "dartsample")
	args = append(args, "--filters.dartsample.radius="+
		strconv.FormatFloat(opts.Radius, 'f', -1, 64))
	args = append(args, "-v", "10", "--debug")
	out, err := exec.Command("pdal", args...).CombinedOutput()

	fmt.Println(string(out))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
