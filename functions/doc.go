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

Crop

Example JSON "options" object for the Crop function.

  {
    "bounds": "([xmin,xmax],[ymin,ymax])",
    "polygon": "POLYGON((30 10, 40 40, 20 40, 10 20, 30 10))",
    "outside": false
  }

Dart

  {
    "radius": 1.0
  }

Dtm

  {
    "grid_size": 1.0
  }

Ground

  {
    "cell_size": 1.0,
    "initial_distance": 0.15,
    "max_distance": 2.5,
    "max_window_size": 33.0,
    "slope": 1.0
  }

Height

Info

  {
    "boundary": false,
    "metadata": true,
    "schema": false
  }

Radius

  {
    "neighbors": 2,
    "radius": 1.0
  }

Statistical

  {
    "neighbors": 8,
    "thresh": 1.5
  }

Translate

  {
    "args": "radiusoutlier ground --filters.radiusoutlier.radius=2.0 --filters.ground.classify=true"
  }

*/
package functions
