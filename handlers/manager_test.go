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
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/venicegeo/pzsvc-pdal/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
)

func TestJobManager(t *testing.T) {
	userJSON := `{
    "status":"submitted"
  }`
	reader := strings.NewReader(userJSON)
	router := httprouter.New()
	router.POST("/manager", JobManagerHandler)
	req, _ := http.NewRequest("POST", "/manager", reader)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	fmt.Println(w.Body)
	if w.Code != http.StatusOK {
		t.Errorf("StatusOK expected: %d", w.Code)
	}
}
