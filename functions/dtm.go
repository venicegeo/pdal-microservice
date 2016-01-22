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

	"github.com/venicegeo/pzsvc-pdal/Godeps/_workspace/src/github.com/venicegeo/pzsvc-sdk-go/job"
)

// DtmOptions defines options for the Dtm function.
type DtmOptions struct {
	GridSize float64 `json:"grid_size"` // size of grid cell in XY dimensions
}

// NewDtmOptions constructs DtmOptions with default values.
func NewDtmOptions() *DtmOptions {
	return &DtmOptions{GridSize: 1.0}
}

// Dtm implements pdal dtm.
func Dtm(w http.ResponseWriter, r *http.Request,
	res *job.OutputMsg, msg job.InputMsg, i, o string) {
	opts := NewDtmOptions()
	if msg.Options != nil {
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			job.BadRequest(w, r, *res, err.Error())
			return
		}
	}

	var args []string
	args = append(args, "translate")
	args = append(args, i)
	args = append(args, "output")
	args = append(args, "ground")
	args = append(args, "--filters.ground.extract=true")
	args = append(args, "--filters.ground.classify=false")
	args = append(args, "-w", "writers.p2g")
	args = append(args, "--writers.p2g.output_type=min")
	args = append(args, "--writers.p2g.output_format=tif")
	args = append(args,
		"--writers.p2g.grid_dist_x="+strconv.FormatFloat(opts.GridSize, 'f', -1, 64))
	args = append(args,
		"--writers.p2g.grid_dist_y="+strconv.FormatFloat(opts.GridSize, 'f', -1, 64))
	args = append(args, "-v10", "--debug")

	out, err := exec.Command("pdal", args...).CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err.Error())
	}
}
