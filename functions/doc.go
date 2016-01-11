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
Package functions provides support for various functions within the PDAL microservice.

Remembering that the PDAL handler expects inputs in the form

	{
		"source": {
			"bucket": <source S3 bucket>,
			"key": <source S3 key>"
		},
		"function": <function name>,
		"options": {
			<function key>: <function value>,
      ...
		},
		"destination": {
			"bucket": <destination S3 bucket>,
			"key" <destination S3 key>
		}
	}

the following sections provide brief examples of valid options for each function. The values provided below do not represent default values.

Crop

Example JSON "options" object for the Crop function.

  {
    "bounds": "([xmin,xmax],[ymin,ymax])",
    "polygon": "POLYGON((30 10, 40 40, 20 40, 10 20, 30 10))",
    "outside": false
  }

Dart

Example JSON "options" object for the Dart function.

  {
    "radius": 1.0
  }

Dtm

Example JSON "options" object for the Dtm function.

  {
    "grid_size": 1.0
  }

Ground

Example JSON "options" object for the Ground function.

  {
    "cell_size": 1.0,
    "initial_distance": 0.15,
    "max_distance": 2.5,
    "max_window_size": 33.0,
    "slope": 1.0
  }

Height

The Height function currently takes no options.

Info

Example JSON "options" object for the Info function.

  {
    "boundary": false,
    "metadata": true,
    "schema": false
  }

Radius

Example JSON "options" object for the Radius function.

  {
    "neighbors": 2,
    "radius": 1.0
  }

Statistical

Example JSON "options" object for the Statistical function.

  {
    "neighbors": 8,
    "thresh": 1.5
  }

Translate

Example JSON "options" object for the Translate function.

  {
    "args": "radiusoutlier ground --filters.radiusoutlier.radius=2.0 --filters.ground.classify=true"
  }

*/
package functions
