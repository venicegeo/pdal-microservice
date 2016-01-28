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

/*
Package utils provides various utilities that are used throughout the project.

Provide functions to return canned responses: StatusOK, StatusBadRequest, and StatusInternalServerError.
*/
package utils

import (
	"log"
	"net/http"
	"os"

	"github.com/venicegeo/pzsvc-sdk-go/job"
	"github.com/venicegeo/pzsvc-sdk-go/s3"
)

// FunctionFunc defines the signature of our function creator.
type FunctionFunc func(
	http.ResponseWriter,
	*http.Request,
	*job.OutputMsg,
	job.InputMsg,
)

func statFile(fname string) {
	log.Printf("Stat %s\n", fname)
	info, _ := os.Stat(fname)
	log.Println(info)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	log.Println("Run!")
}

// MakeFunction wraps the individual PDAL functions.
// Parse the input and output filenames, creating files as needed. Download the
// input data and upload the output data.
func MakeFunction(
	fn func(
		http.ResponseWriter,
		*http.Request,
		*job.OutputMsg,
		job.InputMsg,
		string,
		string,
	),
) FunctionFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		res *job.OutputMsg,
		msg job.InputMsg,
	) {
		log.Println("Here we are in MakeFunction")

		var inputName, outputName string
		var fileIn, fileOut *os.File

		// Split the source S3 key string, interpreting the last element as the
		// input filename. Create the input file, throwing 500 on error.
		inputName = s3.ParseFilenameFromKey(msg.Source.Key)
		fileIn, err := os.Create(inputName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer fileIn.Close()

		// If provided, split the destination S3 key string, interpreting the last
		// element as the output filename. Create the output file, throwing 500 on
		// error.
		if len(msg.Destination.Key) > 0 {
			outputName = s3.ParseFilenameFromKey(msg.Destination.Key)
		}

		// Download the source data from S3, throwing 500 on error.
		err = s3.Download(fileIn, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// statFile(outputName)
		os.Remove(outputName)
		// statFile(outputName)

		// Run the PDAL function.
		fn(w, r, res, msg, inputName, outputName)

		// If an output has been created, upload the destination data to S3,
		// throwing 500 on error.
		if len(msg.Destination.Key) > 0 {
			fileOut, err = os.Open(outputName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer fileOut.Close()
			err = s3.Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
