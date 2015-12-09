# PDAL Microservice

At this point, this repository serves as a sandbox for developing PDAL-based microservices for Piazza.

The going in assumption is that we will receive some message from the dispatcher indicating that a point cloud service has been requested. We will have the path to the data and a description of the task to be performed.

While not the only game in town, PDAL will provide the heavy lifting for most of our point cloud services. We have created a Dockerfile in `venicegeo/dockerfiles/minimal-pdal` that generates a Docker image consisting of PDAL with it's required dependencies and a handful of high-priority plugins (LAZ and NITF support). This in turn serves as the base image for our microservice, which is written in Go.

# Install

```console
$ git clone https://github.com/venicegeo/pdal-microservice
$ ./build.sh
```

The build script will first compile the Go code in a temporary container. The resulting static Go binary is then copied into our `venicegeo/pdal-microservice` image during the `docker build` step. Finally, the service is started on port 8080, mounting your `~/.aws/credentials` to the image.

# Example

This initial example is pretty simple. The service is downloading an LAZ file from our S3 bucket and then calling

```console
$ pdal info <filename> --summary
```

and returning the result. As of this writing, it should look something like

```json
Downloaded file download_file.laz 99563 bytes
{
  "filename": "download_file.laz",
  "pdal_version": "1.1.0 (git-version: 0c36aa)",
  "stats":
  {
    "statistic":
    [
      {
        "average": 512767.0106,
        "count": 38010,
        "maximum": 512834.76,
        "minimum": 512700.87,
        "name": "X",
        "position": 0
      },
      {
        "average": 5403707.591,
        "count": 38010,
        "maximum": 5403849.99,
        "minimum": 5403547.26,
        "name": "Y",
        "position": 1
      },
      {
        "average": 356.1714336,
        "count": 38010,
        "maximum": 404.08,
        "minimum": 295.25,
        "name": "Z",
        "position": 2
      },
      {
        "average": 0.4268350434,
        "count": 38010,
        "maximum": 1,
        "minimum": 0,
        "name": "Intensity",
        "position": 3
      },
      {
        "average": 1,
        "count": 38010,
        "maximum": 1,
        "minimum": 1,
        "name": "ReturnNumber",
        "position": 4
      },
      {
        "average": 1,
        "count": 38010,
        "maximum": 1,
        "minimum": 1,
        "name": "NumberOfReturns",
        "position": 5
      },
      {
        "average": 0,
        "count": 38010,
        "maximum": 0,
        "minimum": 0,
        "name": "ScanDirectionFlag",
        "position": 6
      },
      {
        "average": 0,
        "count": 38010,
        "maximum": 0,
        "minimum": 0,
        "name": "EdgeOfFlightLine",
        "position": 7
      },
      {
        "average": 1.146329913,
        "count": 38010,
        "maximum": 2,
        "minimum": 0,
        "name": "Classification",
        "position": 8
      },
      {
        "average": 0,
        "count": 38010,
        "maximum": 0,
        "minimum": 0,
        "name": "ScanAngleRank",
        "position": 9
      },
      {
        "average": 0,
        "count": 38010,
        "maximum": 0,
        "minimum": 0,
        "name": "UserData",
        "position": 10
      },
      {
        "average": 0,
        "count": 38010,
        "maximum": 0,
        "minimum": 0,
        "name": "PointSourceId",
        "position": 11
      }
    ]
  }
}
```

# Way Forward

Clearly, this is pretty simplistic. We've hardcoded lots of things and have skimped on the error checking. But it demonstrates the capability. Moving forward, we will begin to hash out details of how the PDAL tasks are delivered and make it do some more interesting things, namely executing PDAL pipelines.
