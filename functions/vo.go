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
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/venicegeo/pzsvc-sdk-go/job"
)

// VO implements pdal vo.
func VO(
	w http.ResponseWriter,
	r *http.Request,
	res *job.OutputMsg,
	msg job.InputMsg,
	i, o string,
) {
	log.Println(i)
	log.Println(o)
	out, err := exec.Command("pdal", "vo", i, o, "-v10",
		"--debug").CombinedOutput()

	fmt.Println(string(out))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
