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

package objects

import (
	"encoding/json"
	"time"
)

// StatusType is a string describing the state of the job.
type StatusType int

// Enumerate valid StatusType values.
const (
	Submitted StatusType = iota
	Running
	Success
	Cancelled
	Error
	Fail
)

var statuses = [...]string{"submitted", "running", "success", "cancelled", "error", "fail"}

func (status StatusType) String() string {
	return statuses[status]
}

// S3Bucket defines the expected JSON structure for S3 buckets.
// An S3 bucket can be used for source (input) and destination (output) files.
type S3Bucket struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

// DtmOptions defines options for the DTM generation resource.
type DtmOptions struct {
	GridSize *float64 `json:"grid_size"`
}

// GroundOptions defines options for Ground segmentation.
type GroundOptions struct {
	CellSize        *float64 `json:"cell_size"`
	InitialDistance *float64 `json:"initial_distance"`
	MaxDistance     *float64 `json:"max_distance"`
	MaxWindowSize   *float64 `json:"max_window_size"`
	Slope           *float64 `json:"slope"`
}

// JobInput defines the expected input JSON structure.
// We currently support S3 input (bucket/key), though provider-specific (e.g.,
// GRiD) may be legitimate.
type JobInput struct {
	Source      S3Bucket         `json:"source"`
	Function    *string          `json:"function"`
	Options     *json.RawMessage `json:"options"`
	Destination S3Bucket         `json:"destination"`
}

// JobOutput defines the expected output JSON structure.
type JobOutput struct {
	Input      JobInput                    `json:"input"`
	StartedAt  time.Time                   `json:"started_at"`
	FinishedAt time.Time                   `json:"finished_at"`
	Code       int                         `json:"code"`
	Message    string                      `json:"message"`
	Response   map[string]*json.RawMessage `json:"response"`
}

// JobManagerUpdate defines the expected output JSON structure for updating the JobManager.
type JobManagerUpdate struct {
	Status string `json:"status"`
}
