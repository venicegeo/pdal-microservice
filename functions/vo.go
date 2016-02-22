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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// VO implements pdal vo.
func VO(i, o string, options *json.RawMessage) ([]byte, error) {
	out, err := exec.Command("pdal", "vo", i, o, "-v", "10",
		"--debug").CombinedOutput()

	fmt.Println(string(out))
	if err != nil {
		return nil, err
	}

	fileOut, err := os.Open(o)
	if err != nil {
		return nil, err
	}
	defer fileOut.Close()

	src, err := ioutil.ReadAll(fileOut)
	if err != nil {
		return nil, err
	}

	// Trim whitespace
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, src); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
