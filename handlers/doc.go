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
Package handlers provides support for various handlers within the PDAL microservice.

The primary handler is for calls to the /pdal endpoint. Currently, this parses the S3 bucket/key from the incoming JSON message, downloads the data, and executes the PDAL application with the given function name (e.g., info, translate).

The /pdal resource can be extended by adding additional functions. To add a new function, simply add a case statement such as:

  case "foo":
    utils.MakeFunction(functions.Foo)(w, r, &res, msg)

This will call the MakeFunction routine within the github.com/venicegeo/pzsvc-sdk-go/utils package, downloading source data, processing the data using your custome Foo function, and uploading the result, as needed.

Your custom function should have the following signature.

  type FunctionFunc func(http.ResponseWriter, *http.Request, *job.OutputMsg, job.InputMsg)

The http.ResponseWriter is required so that we can properly communicate errors back. The http.Request was originally included so that we would have access to the original Request - not exactly sure it will be needed. The input and output job messages are required so that we can understand the job that has been dispatched to us and so that we can return the proper response.

We need to add more documentation to the github.com/venicegeo/pzsvc-pdal/functions package, but this is where the magic actually happens. Everything here, in the end, is just a call to PDAL. You should be able to do anything that PDAL can do (depending of course on how you've build PDAL). That includes running kernels (info, translate, merge) and creating pipelines from the CLI. We tend to do more of the latter. Our function options correspond to PDAL CLI arguments. We assume that we always have one input file and one output file (though that's not even a hard requirement). The remainder of the options are parsed and passed. Close examination of any of the existing functions should give you a pretty good sense of what is going on.

The translate function is the all-powerful function (many of the other functions are themselves just translate calls with a little sugar on top). Translate takes a single string as a parameter, which is passed to PDAL. As long as this is something you could run locally with the PDAL CLI, it will run here.

  $ pdal translate <input> <output> [everything else]

For example,

  $ pdal translate <input> <output> ground
  $ pdal translate <input> <output> statisticaloutlier ground height ferry --filters.ferry.dimensions="Height=Z"

Many custom functions will simply be building up these types of PDAL pipelines for common operations. The same readers, writers, and filters can be used in different configurations and with different parameters to produce very different products (e.g., terrain model generation vs. feature extraction).

A secondary handler is a simplified mockup of the Piazza job manager for updating job status.
*/
package handlers
