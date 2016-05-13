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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/venicegeo/pzsvc-pdal/functions"
	"github.com/venicegeo/pzsvc-sdk-go/job"
	"github.com/venicegeo/pzsvc-sdk-go/s3"
)

// InputMsg defines the expected input JSON structure.
// We currently support S3 input (bucket/key), though provider-specific (e.g.,
// GRiD) may be legitimate.
type InputMsg struct {
	Source      interface{}      `json:"source,omitempty"`
	Function    *string          `json:"function,omitempty"`
	Options     *json.RawMessage `json:"options,omitempty"`
	Destination s3.Bucket        `json:"destination,omitempty"`
}

// FunctionFunc defines the signature of our function creator.
type FunctionFunc func(InputMsg) ([]byte, error)

// MakeFunction wraps the individual PDAL functions.
// Parse the input and output filenames, creating files as needed. Download the
// input data and upload the output data.
func MakeFunction(fn func(string, string, *json.RawMessage) ([]byte, error)) FunctionFunc {
	return func(msg InputMsg) ([]byte, error) {
		var inputName, outputName string
		var fileOut *os.File
		log.Printf("%+v\n", msg.Source)
		switch u := msg.Source.(type) {
		case string:
			client := &http.Client{}

			req, err := http.NewRequest("GET", u, nil)
			_, inputName = path.Split(u)

			resp, err := client.Do(req)
			log.Println(resp.Header)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			log.Println(resp.Status)

			fileIn, err := os.Create(inputName)
			if err != nil {
				return nil, err
			}
			defer fileIn.Close()

			numBytes, err := io.Copy(fileIn, resp.Body)
			if err != nil {
				return nil, err
			}

			log.Println("Downloaded", numBytes, "bytes")
		default:
			log.Println("unknown - try as s3 bucket")

			// TODO(chambbj): we know there is a more idiomatic way of achieving this,
			// but it works
			fmt.Printf("%+v\n", msg.Source)
			src := new(s3.Bucket)
			b, err := json.Marshal(msg.Source)
			if err != nil {
				log.Println("error marshaling")
			}
			err = json.Unmarshal(b, &src)
			if err != nil {
				log.Println("must not be an s3 bucket")
			}

			fmt.Println(src.Bucket)
			fmt.Println(src.Key)
			// Split the source S3 key string, interpreting the last element as the
			// input filename. Create the input file, throwing 500 on error.
			inputName = s3.ParseFilenameFromKey(src.Key)
			fileIn, err := os.Create(inputName)
			if err != nil {
				return nil, err
			}
			defer fileIn.Close()

			// Download the source data from S3, throwing 500 on error.
			err = s3.Download(fileIn, src.Bucket, src.Key)
			if err != nil {
				return nil, err
			}
		}

		// If provided, split the destination S3 key string, interpreting the last
		// element as the output filename. Create the output file, throwing 500 on
		// error.
		if len(msg.Destination.Key) > 0 {
			outputName = s3.ParseFilenameFromKey(msg.Destination.Key)
		}

		os.Remove(outputName)

		// Run the PDAL function.
		retval, err := fn(inputName, outputName, msg.Options)
		if err != nil {
			return nil, err
		}

		// If an output has been created, upload the destination data to S3,
		// throwing 500 on error.
		if len(msg.Destination.Key) > 0 {
			fileOut, err = os.Open(outputName)
			if err != nil {
				return nil, err
			}
			defer fileOut.Close()
			err = s3.Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
			if err != nil {
				return nil, err
			}
		}

		return retval, nil
	}
}

// PdalHandler handles PDAL jobs.
func PdalHandler(w http.ResponseWriter, r *http.Request) *AppError {
	// Create the job output message. No matter what happens, we should always be
	// able to populate the StartedAt field.
	var res job.OutputMsg
	res.StartedAt = time.Now()

	var msg InputMsg

	// There should always be a body, else how are we to know what to do? Throw
	// 400 if missing.
	if r.Body == nil {
		return &AppError{nil, "No JSON", http.StatusBadRequest}
	}

	// Throw 500 if we cannot read the body.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	// Throw 400 if we cannot unmarshal the body as a valid InputMsg.
	if err := json.Unmarshal(b, &msg); err != nil {
		return &AppError{err, err.Error(), http.StatusBadRequest}
	}

	// Throw 400 if the JobInput does not specify a function.
	if msg.Function == nil {
		return &AppError{nil, "Must provide a function", http.StatusBadRequest}
	}

	// Make/execute the requested function.
	switch *msg.Function {
	case "crop":
		_, err := MakeFunction(functions.Crop)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "dart":
		_, err := MakeFunction(functions.Dart)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "dtm":
		_, err := MakeFunction(functions.Dtm)(msg)
		if err != nil {
			return &AppError{err, "DTM error", http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "ground":
		_, err := MakeFunction(functions.Ground)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "height":
		_, err := MakeFunction(functions.Height)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "info":
		bytes, err := MakeFunction(functions.Info)(msg)
		if err != nil {
			return &AppError{err, "Info error: MakeFunction", http.StatusInternalServerError}
		}
		if err := json.Unmarshal(bytes, &res.Response); err != nil {
			return &AppError{err, "Info error: Unmarshal", http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "radius":
		_, err := MakeFunction(functions.Radius)(msg)
		if err != nil {
			return &AppError{err, "Radius error", http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "statistical":
		_, err := MakeFunction(functions.Statistical)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "translate":
		_, err := MakeFunction(functions.Translate)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "vo":
		bytes, err := MakeFunction(functions.VO)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(bytes))

	// An unrecognized function will result in 400 error, with message explaining
	// how to list available functions.
	default:
		return &AppError{nil, "Unrecognized function", http.StatusBadRequest}
	}

	return nil
}
