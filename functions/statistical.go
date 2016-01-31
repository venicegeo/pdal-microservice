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
	"strconv"
	"time"

	"github.com/venicegeo/pzsvc-sdk-go/job"
)

// StatisticalOptions defines options for the Statical function.
type StatisticalOptions struct {
	// mean number of neighbors to compute mean and standard deviation
	Neighbors int `json:"neighbors"`
	// standard deviation multiplier for thresholding outliers
	Thresh float64 `json:"thresh"`
}

// NewStatisticalOptions constructs StatisticalOptions with default values.
func NewStatisticalOptions() *StatisticalOptions {
	return &StatisticalOptions{Neighbors: 2, Thresh: 1.5}
}

// Statistical implements pdal height.
func Statistical(
	w http.ResponseWriter,
	r *http.Request,
	res *job.OutputMsg,
	msg job.InputMsg,
	i, o string,
) {
	opts := NewStatisticalOptions()
	if msg.Options != nil {
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var args []string
	args = append(args, "translate", i, o, "statisticaloutlier")
	args = append(args, "--filters.statisticaloutlier.k-neighbors="+
		strconv.Itoa(opts.Neighbors))
	args = append(args, "--filters.statisticaloutlier.thresh="+
		strconv.FormatFloat(opts.Thresh, 'f', -1, 64))
	// we can make this optional later
	args = append(args, "--filters.statisticaloutlier.extract=true")
	args = append(args, "--filters.statisticaloutlier.classify=false")
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
