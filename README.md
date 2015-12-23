# PDAL Microservice

At this point, this repository serves as a sandbox for developing PDAL-based microservices for Piazza.

The going in assumption is that we will receive some message from the dispatcher indicating that a point cloud service has been requested. We will have the path to the data and a description of the task to be performed.

While not the only game in town, PDAL will provide the heavy lifting for most of our point cloud services. We have created a Dockerfile in `venicegeo/dockerfiles/minimal-pdal` that generates a Docker image consisting of PDAL with it's required dependencies and a handful of high-priority plugins (LAZ and NITF support). This in turn serves as the base image for our microservice, which is written in Go.

# Install

```console
$ git clone https://github.com/venicegeo/pdal-microservice
$ scripts/build_and_run.sh
```

The build script will first compile the Go code in a temporary container. The resulting static Go binary is then copied into our `venicegeo/pdal-microservice` image during the `docker build` step. Finally, the service is started on port 8080, mounting your `~/.aws/credentials` to the image.

# Example

Our first example posts the following JSON to the `/pdal` endpoint.

```json
{  
    "source":{  
        "bucket":"venicegeo-sample-data",
        "key":"pointcloud/samp71-utm.laz"
    },
    "function":"info"
}
```

It can be run from the terminal by typing

```console
$ scripts/run-s3-info.sh
```

Internally, the service is simply downloading an LAZ file from our S3 bucket and then calling

```console
$ pdal info <filename>
```

and returning the result. As of this writing, it should look something like

```json
{  
    "input":{  
        "source":{  
            "bucket":"venicegeo-sample-data",
            "key":"pointcloud/samp71-utm.laz"
        },
        "function":"info"
    },
    "started_at":"2015-12-23T18:07:36.987565884Z",
    "finished_at":"2015-12-23T18:07:38.111658707Z",
    "code":200,
    "message":"Success!",
    "response":{  
        "filename":"download_file.laz",
        "pdal_version":"1.1.0 (git-version: 0c36aa)",
        "stats":{  
            "statistic":[  
                {  
                    "average":496348.6372,
                    "count":15645,
                    "maximum":496543.8,
                    "minimum":496148.97,
                    "name":"X",
                    "position":0
                },
                {  
                    "average":5422226.095,
                    "count":15645,
                    "maximum":5422342.88,
                    "minimum":5422121.76,
                    "name":"Y",
                    "position":1
                },
                {  
                    "average":300.0687677,
                    "count":15645,
                    "maximum":309.55,
                    "minimum":293.23,
                    "name":"Z",
                    "position":2
                },
                {  
                    "average":0.113135187,
                    "count":15645,
                    "maximum":1,
                    "minimum":0,
                    "name":"Intensity",
                    "position":3
                },
                {  
                    "average":1,
                    "count":15645,
                    "maximum":1,
                    "minimum":1,
                    "name":"ReturnNumber",
                    "position":4
                },
                {  
                    "average":1,
                    "count":15645,
                    "maximum":1,
                    "minimum":1,
                    "name":"NumberOfReturns",
                    "position":5
                },
                {  
                    "average":0,
                    "count":15645,
                    "maximum":0,
                    "minimum":0,
                    "name":"ScanDirectionFlag",
                    "position":6
                },
                {  
                    "average":0,
                    "count":15645,
                    "maximum":0,
                    "minimum":0,
                    "name":"EdgeOfFlightLine",
                    "position":7
                },
                {  
                    "average":1.773729626,
                    "count":15645,
                    "maximum":2,
                    "minimum":0,
                    "name":"Classification",
                    "position":8
                },
                {  
                    "average":0,
                    "count":15645,
                    "maximum":0,
                    "minimum":0,
                    "name":"ScanAngleRank",
                    "position":9
                },
                {  
                    "average":0,
                    "count":15645,
                    "maximum":0,
                    "minimum":0,
                    "name":"UserData",
                    "position":10
                },
                {  
                    "average":0,
                    "count":15645,
                    "maximum":0,
                    "minimum":0,
                    "name":"PointSourceId",
                    "position":11
                }
            ]
        }
    }
}
```

# Testing

Nothing fancy here. Just run

```console
$ go test ./...
```

Or, if you are interested in code coverage

```console
$ go test ./... -cover
```

And, for more detailed coverage info

```console
$ go test ./... -coverprofile=coverage.out
$ go tool cover -html=coverage.out
```

# Modifying

We use Godeps to aid in deployment. Upon saving, run

```console
$ godep save -r ./...
```

to update the Godeps folder and all import paths.

# Way Forward

Clearly, this is pretty simplistic. We've hardcoded lots of things and have skimped on the error checking. But it demonstrates the capability. Moving forward, we will begin to hash out details of how the PDAL tasks are delivered and make it do some more interesting things, namely executing PDAL pipelines.

# Register Service

```json
{
  "name": "pdal",
  "serviceID": 0,
  "desc": "Process point cloud data using PDAL.",
  "url": "https://api.piazzageo.io/v1/pdal",
  "poc": "",
  "networkAvailable": "TBD",
  "tags": "point cloud, pdal, lidar",
  "classType": "Unclassified",
  "parms": "TBD",
  "termData": "TBD",
  "availability": "UP",
  "serviceQoS": "Development",
  "credentialsRequired": false,
  "clientCert": false,
  "preAuthRequired": false,
  "contracts": ""
}
```
