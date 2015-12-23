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

package objects_test

import (
	"encoding/json"
	"testing"

	"github.com/venicegeo/pdal-microservice/objects"
)

func TestJobInput(t *testing.T) {
	in := `
    {
      "source": {
        "bucket": "Foo",
        "key": "Bar"
      },
      "function": "Baz"
    }`

	b := []byte(in)

	var msg objects.JobInput
	if err := json.Unmarshal(b, &msg); err != nil {
		t.Error("Error parsing JobInput")
	}
	if msg.Source.Bucket != "Foo" {
		t.Error(msg.Source.Bucket, "!= `Foo`")
	}
	if msg.Source.Key != "Bar" {
		t.Error(msg.Source.Key, "!= `Bar`")
	}
	if *msg.Function != "Baz" {
		t.Error(msg.Function, "!= `Baz`")
	}
}

func TestJobOutput(t *testing.T) {
	out := `
    {
      "input": {
        "source": {
          "bucket": "Foo",
          "key": "Bar"
        },
        "function": "Baz"
      },
      "started_at": "2015-12-11T01:31:26.784569058Z",
      "finished_at": "2015-12-11T01:31:26.784569058Z",
      "status": "submitted",
      "response": {"filename":"download_file.laz","pdal_version":"1.1.0 (git-version: 0c36aa)"}
    }`

	b := []byte(out)

	var msg objects.JobOutput
	if err := json.Unmarshal(b, &msg); err != nil {
		t.Error("Error parsing JobOutput")
	}
	if msg.Input.Source.Bucket != "Foo" {
		t.Error(msg.Input.Source.Bucket, "!= `Foo`")
	}
	if msg.Input.Source.Key != "Bar" {
		t.Error(msg.Input.Source.Key, "!= `Bar`")
	}
	if *msg.Input.Function != "Baz" {
		t.Error(msg.Input.Function, "!= `Baz`")
	}
	// if msg.StartedAt != "2015-12-11T01:31:26.784569058Z" {
	// 	t.Error(msg.StartedAt, "!= `2015-12-11T01:31:26.784569058Z`")
	// }
	// if msg.FinishedAt != "2015-12-11T01:31:26.784569058Z" {
	// 	t.Error(msg.FinishedAt, "!= `2015-12-11T01:31:26.784569058Z`")
	// }
	// if msg.Status != "submitted" {
	// 	t.Error(msg.Status, "!= `submitted`")
	// }
	// if msg.Response != `{"filename":"download_file.laz","pdal_version":"1.1.0 (git-version: 0c36aa)"}` {
	// 	t.Error(msg.Response, `!= {"filename":"download_file.laz","pdal_version":"1.1.0 (git-version: 0c36aa)"}`)
	// }
}

func TestStatusTypes(t *testing.T) {
	if objects.Submitted.String() != "submitted" {
		t.Error(objects.Submitted.String(), "!= `submitted`")
	}
	if objects.Running.String() != "running" {
		t.Error(objects.Running.String(), "!= `running`")
	}
	if objects.Success.String() != "success" {
		t.Error(objects.Success.String(), "!= `success`")
	}
	if objects.Cancelled.String() != "cancelled" {
		t.Error(objects.Cancelled.String(), "!= `cancelled`")
	}
	if objects.Error.String() != "error" {
		t.Error(objects.Error.String(), "!= `error`")
	}
	if objects.Fail.String() != "fail" {
		t.Error(objects.Fail.String(), "!= `fail`")
	}
}
