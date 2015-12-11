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

// JobInput defines the expected into JSON structure.
// We currently support S3 input (bucket/key), though provider-specific (e.g.,
// GRiD) may be legitimate.
type JobInput struct {
	Source struct {
		Bucket string `json:"bucket"`
		Key    string `json:"key"`
	} `json:"source"`
	Function string `json:"function"`
}

// JobOutput defines the expected output JSON structure.
type JobOutput struct {
	Input      JobInput                    `json:"input"`
	StartedAt  time.Time                   `json:"started_at"`
	FinishedAt time.Time                   `json:"finished_at"`
	Status     string                      `json:"status"`
	Response   map[string]*json.RawMessage `json:"response"`
}
