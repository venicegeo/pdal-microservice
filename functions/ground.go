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
	"os"
	"os/exec"
	"strconv"

	"github.com/venicegeo/pzsvc-pdal/objects"
	"github.com/venicegeo/pzsvc-pdal/utils"
)

// GroundFunction implements pdal ground.
func GroundFunction(w http.ResponseWriter, r *http.Request,
	res *objects.JobOutput, msg objects.JobInput, f string) {
	fileOut, err := os.Create("output_file.laz")
	if err != nil {
		utils.InternalError(w, r, *res, err.Error())
		return
	}
	defer fileOut.Close()

	cellSize := 1.0
	initialDistance := 0.15
	maxDistance := 2.5
	maxWindowSize := 33.0
	slope := 1.0
	if msg.Options != nil {
		var opts objects.GroundOptions
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			utils.BadRequest(w, r, *res, err.Error())
			return
		}
		if opts.CellSize != nil {
			cellSize = *opts.CellSize
		}
		if opts.InitialDistance != nil {
			initialDistance = *opts.InitialDistance
		}
		if opts.MaxDistance != nil {
			maxDistance = *opts.MaxDistance
		}
		if opts.MaxWindowSize != nil {
			maxWindowSize = *opts.MaxWindowSize
		}
		if opts.Slope != nil {
			slope = *opts.Slope
		}
	}
	cellSizeStr := "--filters.ground.cell_size=" +
		strconv.FormatFloat(cellSize, 'f', -1, 64)
	initialDistanceStr := "--filters.ground.initial_distance=" +
		strconv.FormatFloat(initialDistance, 'f', -1, 64)
	maxDistanceStr := "--filters.ground.max_distance=" +
		strconv.FormatFloat(maxDistance, 'f', -1, 64)
	maxWindowSizeStr := "--filters.ground.max_window_size=" +
		strconv.FormatFloat(maxWindowSize, 'f', -1, 64)
	slopeStr := "--filters.ground.slope=" +
		strconv.FormatFloat(slope, 'f', -1, 64)
	fmt.Println(maxWindowSizeStr)

	out, err := exec.Command("pdal", "translate", f, fileOut.Name(),
		"ground", "--filters.ground.extract=true",
		"--filters.ground.classify=false", cellSizeStr, initialDistanceStr,
		maxDistanceStr, maxWindowSizeStr, slopeStr, "-v10",
		"--debug").CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err.Error())
	}

	err = utils.S3Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
	if err != nil {
		utils.InternalError(w, r, *res, err.Error())
		return
	}
}
