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
	"os/exec"
	"strconv"
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
func Statistical(i, o string, options *json.RawMessage) ([]byte, error) {
	opts := NewStatisticalOptions()
	if options != nil {
		if err := json.Unmarshal(*options, &opts); err != nil {
			return nil, err
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
		return nil, err
	}

	return nil, nil
}
