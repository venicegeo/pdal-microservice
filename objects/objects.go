package objects

import (
	"encoding/json"
	"time"
)

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
