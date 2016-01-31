/*
Copyright 2016, RadiantBlue Technologies, Inc.

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
	"time"

	"github.com/venicegeo/pzsvc-sdk-go/job"
)

// CropOptions defines options for the Crop function.
type CropOptions struct {
	// extents of the clipping rectangle in the form "([xmin,xmax],[ymin,ymax])"
	Bounds string `json:"bounds"`
	// the clipping polygon in well-known text, e.g.,
	// POLYGON((30 10, 40 40, 20 40, 10 20, 30 10))
	Polygon string `json:"polygon"`
	// invert logic and only keep points outside the bounds/polygon
	// (default: false)
	Outside bool `json:"outside"`
}

// NewCropOptions constructs CropOptions with default values.
func NewCropOptions() *CropOptions {
	return &CropOptions{Outside: false}
}

/*
Crop calls PDAL translate with a crop filter.

The Crop function will invoke the PDAL translate command as follows:

	$ pdal translate <input> <output> crop \
	  [--filters.crop.bounds=<bounds string>] \
	  [--filters.crop.polygon=<polygon string>] \
	  [--filters.crop.outside=<true|false>] \
	  -v10 --debug
*/
func Crop(
	w http.ResponseWriter,
	r *http.Request,
	res *job.OutputMsg,
	msg job.InputMsg,
	i, o string,
) {
	opts := NewCropOptions()
	if msg.Options != nil {
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var args []string
	args = append(args, "translate", i, o, "crop")
	if (opts.Bounds == "" && opts.Polygon == "") ||
		(opts.Bounds != "" && opts.Polygon != "") {
		fmt.Println("must provide bounds OR polygon, but not both")
	}
	if opts.Bounds != "" {
		args = append(args, "--filters.crop.bounds="+opts.Bounds)
	} else if opts.Polygon != "" {
		args = append(args, "--filters.crop.polygon="+opts.Polygon)
	}
	if opts.Outside {
		args = append(args, "--filters.crop.outside=true")
	} else {
		args = append(args, "--filters.crop.outside=false")
	}
	args = append(args, "-v10", "--debug")
	out, err := exec.Command("pdal", args...).CombinedOutput()

	fmt.Println(string(out))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If we made it here, we can record the FinishedAt time, notify the job
	// manager of success, and return 200.
	res.FinishedAt = time.Now()
	job.Okay(w, r, *res, "Success!")
}
