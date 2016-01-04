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

package functions

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/venicegeo/pzsvc-pdal/objects"
	"github.com/venicegeo/pzsvc-pdal/utils"
)

// HeightFunction implements pdal height.
func HeightFunction(w http.ResponseWriter, r *http.Request,
	res *objects.JobOutput, msg objects.JobInput, f string) {
	fileOut, err := os.Create("output_file.laz")
	if err != nil {
		utils.InternalError(w, r, *res, err.Error())
		return
	}
	defer fileOut.Close()

	out, err := exec.Command("pdal", "translate", f, fileOut.Name(),
		"ground", "height", "ferry",
		"--filters.ferry.dimensions=Height=Z", "-v10", "--debug").CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		utils.InternalError(w, r, *res, err.Error())
		return
	}

	err = utils.S3Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
	if err != nil {
		utils.InternalError(w, r, *res, err.Error())
		return
	}
}
