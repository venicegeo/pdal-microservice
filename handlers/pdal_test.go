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

// these require PDAL - not appropriate for Jenkins
// func TestBasicInfo(t *testing.T) {
// 	userJSON := `{
// 		"source":
// 		{
// 			"bucket": "venicegeo-sample-data",
// 			"key": "pointcloud/samp71-utm.laz"
// 		},
// 		"function": "info"
// 	}`
// 	reader := strings.NewReader(userJSON)
// 	router := httprouter.New()
// 	router.POST("/pdal", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/pdal", reader)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusOK {
// 		t.Errorf("StatusOK expected: %d", w.Code)
// 	}
// }
//
// func TestBasicGround(t *testing.T) {
// 	userJSON := `{
// 		"source":
// 		{
// 			"bucket": "venicegeo-sample-data",
// 			"key": "pointcloud/samp71-utm.laz"
// 		},
// 		"function": "ground",
// 		"destination":
// 		{
// 			"bucket": "venicegeo-sample-data",
// 			"key": "temp/ground.laz"
// 		}
// 	}`
// 	reader := strings.NewReader(userJSON)
// 	router := httprouter.New()
// 	router.POST("/pdal", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/pdal", reader)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusOK {
// 		t.Errorf("StatusOK expected: %d", w.Code)
// 	}
// }
//
// func TestBasicHeight(t *testing.T) {
// 	userJSON := `{
// 		"source":
// 		{
// 			"bucket": "venicegeo-sample-data",
// 			"key": "pointcloud/samp71-utm.laz"
// 		},
// 		"function": "height",
// 		"destination": {
// 			"bucket": "venicegeo-sample-data",
// 			"key": "temp/height.laz"
// 		}
// 	}`
// 	reader := strings.NewReader(userJSON)
// 	router := httprouter.New()
// 	router.POST("/pdal", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/pdal", reader)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusOK {
// 		t.Errorf("StatusOK expected: %d", w.Code)
// 	}
// }

// // this is specific to our service, it stays
// func TestNoFunctionField(t *testing.T) {
// 	userJSON := `{
// 		"source":
// 		{
// 			"bucket": "venicegeo-sample-data",
// 			"key": "pointcloud/samp71-utm.laz"
// 		},
// 		"fail": "info"
// 	}`
// 	reader := strings.NewReader(userJSON)
// 	router := httprouter.New()
// 	router.POST("/api/v1/pdal", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/api/v1/pdal", reader)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusBadRequest {
// 		t.Errorf("StatusBadRequest expected: %d", w.Code)
// 	}
// }
//
// // also specific to our service, it too stays
// func TestBadFunction(t *testing.T) {
// 	userJSON := `{
// 		"source":
// 		{
// 			"bucket": "venicegeo-sample-data",
// 			"key": "pointcloud/samp71-utm.laz"
// 		},
// 		"function": "fail"
// 	}`
// 	reader := strings.NewReader(userJSON)
// 	router := httprouter.New()
// 	router.POST("/api/v1/pdal", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/api/v1/pdal", reader)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusBadRequest {
// 		t.Errorf("StatusBadRequest expected: %d", w.Code)
// 	}
// }

// these really test elements of pzsvc-sdk-go
// func TestBadBucket(t *testing.T) {
// 	userJSON := `{
// 		"source":
// 		{
// 			"bucket": "bad-bucket",
// 			"key": "pointcloud/samp71-utm.laz"
// 		},
// 		"function": "info"
// 	}`
// 	reader := strings.NewReader(userJSON)
// 	router := httprouter.New()
// 	router.POST("/pdal", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/pdal", reader)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusInternalServerError {
// 		t.Errorf("StatusInternalServerError expected: %d", w.Code)
// 	}
// }
//
// func TestBadKey(t *testing.T) {
// 	userJSON := `{
// 		"source":
// 		{
// 			"bucket": "venicegeo-sample-data",
// 			"key": "bad-folder/bad-file"
// 		},
// 		"function": "info"
// 	}`
// 	reader := strings.NewReader(userJSON)
// 	router := httprouter.New()
// 	router.POST("/pdal", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/pdal", reader)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusInternalServerError {
// 		t.Errorf("StatusInternalServerError expected: %d", w.Code)
// 	}
// }
//
// func TestEmptyJSON(t *testing.T) {
// 	userJSON := `{}`
// 	reader := strings.NewReader(userJSON)
// 	router := httprouter.New()
// 	router.POST("/pdal", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/pdal", reader)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusBadRequest {
// 		t.Errorf("StatusBadRequest expected: %d", w.Code)
// 	}
// }
//
// func TestNoJSON(t *testing.T) {
// 	router := httprouter.New()
// 	router.POST("/pdal", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/pdal", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusBadRequest {
// 		t.Errorf("StatusBadRequest expected: %d", w.Code)
// 	}
// }
//
// func TestBadEndpoint(t *testing.T) {
// 	router := httprouter.New()
// 	router.POST("/ladp", PdalHandler)
// 	req, _ := http.NewRequest("POST", "/ladp", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	fmt.Println(w.Body)
// 	if w.Code != http.StatusBadRequest {
// 		t.Errorf("StatusBadRequest expected: %d", w.Code)
// 	}
// }
