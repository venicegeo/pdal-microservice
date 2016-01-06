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

	"github.com/venicegeo/pzsvc-pdal/Godeps/_workspace/src/github.com/venicegeo/pzsvc-sdk-go/objects"
	"github.com/venicegeo/pzsvc-pdal/Godeps/_workspace/src/github.com/venicegeo/pzsvc-sdk-go/utils"
)

// GroundOptions defines options for Ground segmentation.
type GroundOptions struct {
	CellSize        float64 `json:"cell_size"`
	InitialDistance float64 `json:"initial_distance"`
	MaxDistance     float64 `json:"max_distance"`
	MaxWindowSize   float64 `json:"max_window_size"`
	Slope           float64 `json:"slope"`
}

// NewGroundOptions constructs GroundOptions with default values.
func NewGroundOptions() *GroundOptions {
	return &GroundOptions{
		CellSize:        1.0,
		InitialDistance: 0.15,
		MaxDistance:     2.5,
		MaxWindowSize:   33.0,
		Slope:           1.0,
	}
}

// GroundFunction implements pdal ground.
func GroundFunction(w http.ResponseWriter, r *http.Request,
	res *objects.JobOutput, msg objects.JobInput, i, o string) {
	opts := NewGroundOptions()
	if msg.Options != nil {
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			utils.BadRequest(w, r, *res, err.Error())
			return
		}
	}

	var args []string
	args = append(args, "translate")
	args = append(args, i)
	args = append(args, o)
	args = append(args, "ground")
	args = append(args, "--filters.ground.extract=true")
	args = append(args, "--filters.ground.classify=false")
	args = append(args,
		"--filters.ground.cell_size="+strconv.FormatFloat(opts.CellSize, 'f', -1, 64))
	args = append(args,
		"--filters.ground.initial_distance="+strconv.FormatFloat(opts.InitialDistance, 'f', -1, 64))
	args = append(args,
		"--filters.ground.max_distance="+strconv.FormatFloat(opts.MaxDistance, 'f', -1, 64))
	args = append(args,
		"--filters.ground.max_window_size="+strconv.FormatFloat(opts.MaxWindowSize, 'f', -1, 64))
	args = append(args,
		"--filters.ground.slope="+strconv.FormatFloat(opts.Slope, 'f', -1, 64))
	args = append(args, "-v10", "--debug")

	out, err := exec.Command("pdal", args...).CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err.Error())
	}
}
