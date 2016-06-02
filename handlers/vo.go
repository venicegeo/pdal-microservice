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

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
)

// VoOptions defines options for the VO function.
type VoOptions struct {
	Filename    *string `json:"filename"` // list of input file paths
	AGL         float64 `json:"agl"`
	MaxSize     uint16  `json:"max_size"`
	MinSize     uint16  `json:"min_size"`
	Resolution  float64 `json:"resolution"`
	Tolerance   float64 `json:"tolerance"`
	Z0Tolerance float64 `json:"z0_tolerance"`
	Denoise     bool    `json:"denoise"`
	AssignSRS   *string `json:"a_srs"`
}

// NewVoOptions constructs VoOptions with default values.
func NewVoOptions() *VoOptions {
	return &VoOptions{AGL: 20.0, MaxSize: 65535, MinSize: 1, Resolution: 10, Tolerance: 3, Z0Tolerance: 6.0}
}

func getFloatAsString(name string, val float64) string {
	return name + "=" + strconv.FormatFloat(val, 'f', -1, 64)
}

func getUintAsString(name string, val uint16) string {
	return name + "=" + strconv.FormatUint(uint64(val), 10)
}

// VoHandler handles PDAL jobs.
func VoHandler(w http.ResponseWriter, r *http.Request) *AppError {
	if r.Body == nil {
		return &AppError{nil, "No JSON", http.StatusBadRequest}
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	opts := NewVoOptions()
	if b != nil {
		if err := json.Unmarshal(b, &opts); err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}
	}

	if opts.Filename == nil {
		return &AppError{err, "No filename", http.StatusInternalServerError}
	}

	s3url, err := url.Parse(*opts.Filename)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	name, _, err := inferAndDownload(s3url)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	outFile := "out.json"
	os.Remove(outFile)

	if opts.AssignSRS != nil {
		projected := name + "-projected.laz"

		projStr := "--writers.las.a_srs=" + *opts.AssignSRS

		srsArgs := []string{
			"translate", name, projected, projStr, "-v", "3", "--debug",
		}
		log.Println("Assigning SRS with args", srsArgs)
		srsOut, srsErr := exec.Command("pdal", srsArgs...).CombinedOutput()
		if err != nil {
			return &AppError{srsErr, srsErr.Error(), http.StatusInternalServerError}
		}
		log.Println("PDAL CLI responded with")
		log.Println(string(srsOut))

		name = projected
	}

	args := []string{
		"translate", name, outFile,
		"-w", "writers.vo", "-v", "3", "--debug",
		getFloatAsString("--writers.vo.agl", opts.AGL),
		getUintAsString("--writers.vo.max_size", opts.MaxSize),
		getUintAsString("--writers.vo.min_size", opts.MinSize),
		getFloatAsString("--writers.vo.resolution", opts.Resolution),
		getFloatAsString("--writers.vo.tolerance", opts.Tolerance),
		getFloatAsString("--writers.vo.z0_tolerance", opts.Z0Tolerance),
	}
	if opts.Denoise {
		args = append(args, "-f", "filters.statisticaloutlier", "--filters.statisticaloutlier.extract=true")
	}
	log.Println("PDAL CLI called with args", args)
	out, err := exec.Command("pdal", args...).CombinedOutput()
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}
	log.Println("PDAL CLI responded with")
	log.Println(string(out))

	fileOut, err := os.Open(outFile)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}
	defer fileOut.Close()

	src, err := ioutil.ReadAll(fileOut)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, src); err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, buffer.String())

	// TODO(chambbj): there is a case where we have assigned a projected, which
	// created another file that needs to be cleaned up...may be better in the
	// future to create/use a temp dir that just gets cleared, wiping everything
	os.Remove(name)

	return nil
}
