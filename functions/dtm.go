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

package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/venicegeo/pzsvc-pdal/objects"
	"github.com/venicegeo/pzsvc-pdal/utils"
)

// DtmFunction implements pdal dtm.
func DtmFunction(w http.ResponseWriter, r *http.Request,
	res *objects.JobOutput, msg objects.JobInput, i, o string) {
	gridSize := 1.0
	if msg.Options != nil {
		var opts objects.DtmOptions
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			utils.BadRequest(w, r, *res, err.Error())
			return
		}
		if opts.GridSize != nil {
			gridSize = *opts.GridSize
		}
	}
	fmt.Println(gridSize)
	gridDistX := "--writers.p2g.grid_dist_x=" +
		strconv.FormatFloat(gridSize, 'f', -1, 64)
	gridDistY := "--writers.p2g.grid_dist_y=" +
		strconv.FormatFloat(gridSize, 'f', -1, 64)

	out, err := exec.Command("pdal", "translate", i, "output",
		"ground", "--filters.ground.extract=true",
		"--filters.ground.classify=false", "-w", "writers.p2g",
		"--writers.p2g.output_type=min", "--writers.p2g.output_format=tif",
		gridDistX, gridDistY, "-v10", "--debug").CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err.Error())
	}
}
