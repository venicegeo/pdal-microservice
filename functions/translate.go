/*
Copyright 2015-2016, RadiantBlue Technologies, Inc.

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
	"strings"

	"github.com/venicegeo/pzsvc-sdk-go/job"
)

// TranslateOptions defines options for Ground segmentation.
type TranslateOptions struct {
	Args string `json:"args"`
}

// NewTranslateOptions constructs TranslateOptions with default values.
func NewTranslateOptions() *TranslateOptions {
	return &TranslateOptions{Args: ""}
}

// Translate implements pdal translate.
func Translate(w http.ResponseWriter, r *http.Request,
	res *job.OutputMsg, msg job.InputMsg, i, o string) {
	opts := NewTranslateOptions()
	if msg.Options != nil {
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			job.BadRequest(w, r, *res, err.Error())
			return
		}
	}

	optArgs := strings.Split(opts.Args, " ")

	var args []string
	args = append(args, "translate")
	args = append(args, i)
	args = append(args, o)
	args = append(args, optArgs...)
	args = append(args, "-v10", "--debug")

	out, err := exec.Command("pdal", args...).CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err.Error())
	}
}
