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

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/venicegeo/pzsvc-pdal/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-pdal/functions"
	"github.com/venicegeo/pzsvc-pdal/objects"
	"github.com/venicegeo/pzsvc-pdal/utils"
)

type functionFunc func(http.ResponseWriter, *http.Request,
	*objects.JobOutput, objects.JobInput)

// makeFunction wraps the individual PDAL functions.
// Parse the input and output filenames, creating files as needed. Download the
// input data and upload the output data.
func makeFunction(fn func(http.ResponseWriter, *http.Request,
	*objects.JobOutput, objects.JobInput, string, string)) functionFunc {
	return func(w http.ResponseWriter, r *http.Request, res *objects.JobOutput,
		msg objects.JobInput) {
		var inputName, outputName string
		var fileIn, fileOut *os.File

		// Split the source S3 key string, interpreting the last element as the
		// input filename. Create the input file, throwing 500 on error.
		keySlice := strings.Split(msg.Source.Key, "/")
		inputName = keySlice[len(keySlice)-1]
		fileIn, err := os.Create(inputName)
		if err != nil {
			utils.InternalError(w, r, *res, err.Error())
			return
		}
		defer fileIn.Close()

		// If provided, split the destination S3 key string, interpreting the last
		// element as the output filename. Create the output file, throwing 500 on
		// error.
		if len(msg.Destination.Key) > 0 {
			keySlice = strings.Split(msg.Destination.Key, "/")
			outputName = keySlice[len(keySlice)-1]
			fileOut, err = os.Create(outputName)
			if err != nil {
				utils.InternalError(w, r, *res, err.Error())
				return
			}
			defer fileOut.Close()
		}

		// Download the source data from S3, throwing 500 on error.
		err = utils.S3Download(fileIn, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			utils.InternalError(w, r, *res, err.Error())
			return
		}

		// Run the PDAL function.
		fn(w, r, res, msg, inputName, outputName)

		// If an output has been created, upload the destination data to S3,
		// throwing 500 on error.
		if len(msg.Destination.Key) > 0 {
			err = utils.S3Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
			if err != nil {
				utils.InternalError(w, r, *res, err.Error())
				return
			}
		}
	}
}

// PdalHandler handles PDAL jobs.
func PdalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Create the job output message. No matter what happens, we should always be
	// able to populate the StartedAt field.
	var res objects.JobOutput
	res.StartedAt = time.Now()

	// There should always be a body, else how are we to know what to do? Throw
	// 400 if missing.
	if r.Body == nil {
		utils.BadRequest(w, r, res, "No JSON")
		return
	}

	// Throw 500 if we cannot read the body.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.InternalError(w, r, res, err.Error())
		return
	}

	// Throw 400 if we cannot unmarshal the body as a valid JobInput.
	var msg objects.JobInput
	if err := json.Unmarshal(b, &msg); err != nil {
		utils.BadRequest(w, r, res, err.Error())
		return
	}

	// Throw 400 if the JobInput does not specify a function.
	if msg.Function == nil {
		utils.BadRequest(w, r, res, "Must provide a function")
		return
	}

	// If everything is okay up to this point, we will echo the JobInput in the
	// JobOutput and mark the job as Running.
	res.Input = msg
	utils.UpdateJobManager(objects.Running, r)

	// Make/execute the requested function.
	switch *msg.Function {
	case "info":
		makeFunction(functions.InfoFunction)(w, r, &res, msg)

	case "ground":
		makeFunction(functions.GroundFunction)(w, r, &res, msg)

	case "height":
		makeFunction(functions.HeightFunction)(w, r, &res, msg)

	case "dtm":
		makeFunction(functions.DtmFunction)(w, r, &res, msg)

	case "dart":
		makeFunction(functions.DartFunction)(w, r, &res, msg)

	case "list":
		out := []byte(`{"functions":["info","ground","height","dtm","dart","list"]}`)

		if err := json.Unmarshal(out, &res.Response); err != nil {
			log.Fatal(err)
		}

	// list options for named function
	case "options":
		// start with only info options as a test
		foo := functions.NewInfoOptions()
		bar, _ := json.Marshal(foo)
		if err := json.Unmarshal(bar, &res.Response); err != nil {
			log.Fatal(err)
		}

	// An unrecognized function will result in 400 error, with message explaining
	// how to list available functions.
	default:
		utils.BadRequest(w, r, res, "")
		return
	}

	// If we made it here, we can record the FinishedAt time, notify the job
	// manager of success, and return 200.
	res.FinishedAt = time.Now()
	utils.Okay(w, r, res, "Success!")
}
