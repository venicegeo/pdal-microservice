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

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pdal-microservice/objects"
	"github.com/venicegeo/pdal-microservice/utils"
)

// var validPath = regexp.MustCompile("^/(info|pipeline)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(pdal)$")

// PdalHandler handles PDAL jobs.
func PdalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var res objects.JobOutput
	res.StartedAt = time.Now()

	// Check that we have a valid path. Is this the correct place to do this?
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		utils.BadRequest(w, r, res, "Endpoint does not exist")
		return
	}

	if r.Body == nil {
		utils.BadRequest(w, r, res, "No JSON")
		return
	}

	// Parse the incoming JSON body, and unmarshal as events.NewData struct.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.InternalError(w, r, res, err.Error())
		return
	}

	var msg objects.JobInput
	if err := json.Unmarshal(b, &msg); err != nil {
		utils.BadRequest(w, r, res, err.Error())
		return
	}
	if msg.Function == nil {
		utils.BadRequest(w, r, res, "Must provide a function")
		return
	}

	res.Input = msg
	utils.UpdateJobManager(objects.Running, r)

	switch *msg.Function {
	case "info":
		file, err := os.Create("download_file.laz")
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}
		defer file.Close()

		err = utils.S3Download(file, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			utils.InternalError(w, r, res, err.Error())
			return
		}

		out, _ := exec.Command("pdal", *msg.Function, file.Name()).CombinedOutput()

		// Trim whitespace
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, out); err != nil {
			fmt.Println(err)
		}

		if err = json.Unmarshal(buffer.Bytes(), &res.Response); err != nil {
			log.Fatal(err)
		}

	case "pipeline":
		fmt.Println("pipeline not implemented yet")

	default:
		utils.BadRequest(w, r, res,
			"Only the info and pipeline functions are supported at this time")
		return
	}

	res.FinishedAt = time.Now()
	utils.Okay(w, r, res, "Success!")
}
