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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/venicegeo/pzsvc-pdal/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-pdal/functions"
	"github.com/venicegeo/pzsvc-pdal/objects"
	"github.com/venicegeo/pzsvc-pdal/utils"
)

type functionFunc func(http.ResponseWriter, *http.Request,
	*objects.JobOutput, objects.JobInput)

func makeIFunction(fn func(http.ResponseWriter, *http.Request,
	*objects.JobOutput, objects.JobInput, string)) functionFunc {
	return func(w http.ResponseWriter, r *http.Request, res *objects.JobOutput,
		msg objects.JobInput) {
		file, err := os.Create("download_file.laz")
		if err != nil {
			utils.InternalError(w, r, *res, err.Error())
			return
		}
		defer file.Close()

		err = utils.S3Download(file, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			utils.InternalError(w, r, *res, err.Error())
			return
		}
		fn(w, r, res, msg, file.Name())
	}
}

func makeIOFunction(fn func(http.ResponseWriter, *http.Request,
	*objects.JobOutput, objects.JobInput, string, string)) functionFunc {
	return func(w http.ResponseWriter, r *http.Request, res *objects.JobOutput,
		msg objects.JobInput) {
		file, err := os.Create("download_file.laz")
		if err != nil {
			utils.InternalError(w, r, *res, err.Error())
			return
		}
		defer file.Close()

		fileOut, err := os.Create("output.min.tif")
		if err != nil {
			utils.InternalError(w, r, *res, err.Error())
			return
		}
		defer fileOut.Close()

		err = utils.S3Download(file, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			utils.InternalError(w, r, *res, err.Error())
			return
		}

		fn(w, r, res, msg, file.Name(), fileOut.Name())

		err = utils.S3Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
		if err != nil {
			utils.InternalError(w, r, *res, err.Error())
			return
		}
	}
}

// PdalHandler handles PDAL jobs.
func PdalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var res objects.JobOutput
	res.StartedAt = time.Now()

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
		makeIFunction(functions.InfoFunction)(w, r, &res, msg)

	case "ground":
		makeIOFunction(functions.GroundFunction)(w, r, &res, msg)

	case "height":
		makeIOFunction(functions.HeightFunction)(w, r, &res, msg)

	case "dtm":
		makeIOFunction(functions.DtmFunction)(w, r, &res, msg)

	// list available functions

	// list options for named function

	default:
		utils.BadRequest(w, r, res, "Send message telling user to pass 'list' as the function to see a list of available functions.")
		return
	}

	res.FinishedAt = time.Now()
	utils.Okay(w, r, res, "Success!")
}
