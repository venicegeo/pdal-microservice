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

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/venicegeo/pzsvc-pdal/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-pdal/objects"
	"github.com/venicegeo/pzsvc-pdal/utils"
)

// PdalHandler handles PDAL jobs.
func PdalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var res objects.JobOutput
	res.StartedAt = time.Now()

	if r.Body == nil {
		utils.BadRequest(w, r, res, "No JSON")
		return
	}

	// Parse the incoming JSON body, and unmarshal as events.NewData struct.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.InternalError(w, r, res, err.Error())
		return
	}

	var msg objects.JobInput
	if err := json.Unmarshal(b, &msg); err != nil {
		utils.BadRequest(w, r, res, err.Error())
		return
	}
	if msg.Function == nil {
		utils.BadRequest(w, r, res, "Must provide a function")
		return
	}

	res.Input = msg
	utils.UpdateJobManager(objects.Running, r)

	switch *msg.Function {
	case "info":
		file, err := os.Create("download_file.laz")
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}
		defer file.Close()

		err = utils.S3Download(file, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}

		out, _ := exec.Command("pdal", *msg.Function, file.Name()).CombinedOutput()

		// Trim whitespace
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, out); err != nil {
			fmt.Println(err)
		}

		if err = json.Unmarshal(buffer.Bytes(), &res.Response); err != nil {
			log.Fatal(err)
		}

	case "pipeline":
		fmt.Println("pipeline not implemented yet")

	case "ground":
		file, err := os.Create("download_file.laz")
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}
		defer file.Close()

		fileOut, err := os.Create("output_file.laz")
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}
		defer fileOut.Close()

		err = utils.S3Download(file, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}

		cellSize := 1.0
		initialDistance := 0.15
		maxDistance := 2.5
		maxWindowSize := 33.0
		slope := 1.0
		if msg.Options != nil {
			var opts objects.GroundOptions
			if err := json.Unmarshal(*msg.Options, &opts); err != nil {
				utils.BadRequest(w, r, res, err.Error())
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

		out, err := exec.Command("pdal", "translate", file.Name(), fileOut.Name(),
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
			utils.InternalError(w, r, res, err.Error())
			return
		}

	case "height":
		file, err := os.Create("download_file.laz")
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}
		defer file.Close()

		fileOut, err := os.Create("output_file.laz")
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}
		defer fileOut.Close()

		err = utils.S3Download(file, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}

		out, err := exec.Command("pdal", "translate", file.Name(), fileOut.Name(),
			"ground", "height", "ferry",
			"--filters.ferry.dimensions=Height=Z", "-v10", "--debug").CombinedOutput()

		if err != nil {
			fmt.Println(string(out))
			utils.InternalError(w, r, res, err.Error())
			return
		}

		err = utils.S3Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}

	case "groundopts":
		_, err := exec.Command("pdal",
			"--options=filters.ground").CombinedOutput()

		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}

	case "dtm":
		file, err := os.Create("download_file.laz")
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}
		defer file.Close()

		fileOut, err := os.Create("output.min.tif")
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}
		defer fileOut.Close()

		err = utils.S3Download(file, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}

		gridSize := 1.0
		if msg.Options != nil {
			var opts objects.DtmOptions
			if err := json.Unmarshal(*msg.Options, &opts); err != nil {
				utils.BadRequest(w, r, res, err.Error())
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

		out, err := exec.Command("pdal", "translate", file.Name(), "output",
			"ground", "--filters.ground.extract=true",
			"--filters.ground.classify=false", "-w", "writers.p2g",
			"--writers.p2g.output_type=min", "--writers.p2g.output_format=tif",
			gridDistX, gridDistY, "-v10", "--debug").CombinedOutput()

		if err != nil {
			fmt.Println(string(out))
			fmt.Println(err.Error())
		}

		err = utils.S3Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}

	/*
		I get a bad_alloc here, but only via go test. The same command run natively
		works fine.
		case "drivers":
			out, err := exec.Command("pdal",
				"--drivers").CombinedOutput()

			fmt.Println(string(out))
			if err != nil {
				utils.InternalError(w, r, res, err.Error())
				return
			}
	*/

	default:
		utils.BadRequest(w, r, res,
			"Only the info and pipeline functions are supported at this time")
		return
	}

	res.FinishedAt = time.Now()
	utils.Okay(w, r, res, "Success!")
}
