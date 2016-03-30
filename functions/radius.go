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

// RadiusOptions defines options for the Radius function.
type RadiusOptions struct {
	// minimum number of neighbors in radius
	Neighbors int     `json:"neighbors"`
	Radius    float64 `json:"radius"`
}

// NewRadiusOptions constructs RadiusOptions with default values.
func NewRadiusOptions() *RadiusOptions {
	return &RadiusOptions{Neighbors: 2, Radius: 1.0}
}

// Radius implements pdal height.
func Radius(i, o string, options *json.RawMessage) ([]byte, error) {
	opts := NewRadiusOptions()
	if options != nil {
		if err := json.Unmarshal(*options, &opts); err != nil {
			return nil, err
		}
	}

	var args []string
	args = append(args, "translate", i, o, "radiusoutlier")
	args = append(args, "--filters.radiusoutlier.min_neighbors="+
		strconv.Itoa(opts.Neighbors))
	args = append(args, "--filters.radiusoutlier.radius="+
		strconv.FormatFloat(opts.Radius, 'f', -1, 64))
	// we can make this optional later
	args = append(args, "--filters.radiusoutlier.extract=true")
	args = append(args, "--filters.radiusoutlier.classify=false")
	args = append(args, "-v", "10", "--debug")
	out, err := exec.Command("pdal", args...).CombinedOutput()

	fmt.Println(string(out))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
