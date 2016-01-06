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
	"net/http"
	"os/exec"
	"strconv"

	"github.com/venicegeo/pzsvc-sdk-go/job"
)

// DartOptions defines options for dart sampling.
type DartOptions struct {
	// Radius is minimum distance criteria. No two points in the sampled point
	// cloud will be closer than the specified radius.
	Radius float64 `json:"radius"`
}

// NewDartOptions constructs DartOptions with default values.
func NewDartOptions() *DartOptions {
	return &DartOptions{Radius: 1.0}
}

// DartFunction implements pdal height.
func DartFunction(w http.ResponseWriter, r *http.Request,
	res *job.OutputMsg, msg job.InputMsg, i, o string) {
	opts := NewDartOptions()
	if msg.Options != nil {
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			job.BadRequest(w, r, *res, err.Error())
			return
		}
	}

	var args []string
	args = append(args, "translate", i, o, "dartsample")
	args = append(args,
		"--filters.dartsample.radius="+strconv.FormatFloat(opts.Radius, 'f', -1, 64))
	args = append(args, "-v10", "--debug")
	out, err := exec.Command("pdal", args...).CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err.Error())
	}
}
