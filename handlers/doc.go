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

/*
Package handlers provides support for various handlers within the PDAL microservice.

The primary handler is for calls to the /pdal endpoint. Currently, this parses the S3 bucket/key from the incoming JSON message, downloads the data, and executes the PDAL application with the given function name (e.g., info, translate).

A secondary handler is a simplified mockup of the Piazza job manager for updating job status.
*/
package handlers
