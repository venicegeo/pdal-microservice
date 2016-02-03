[![GoDoc](https://godoc.org/github.com/venicegeo/pzsvc-pdal?status.svg)](https://godoc.org/github.com/venicegeo/pzsvc-pdal)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/venicegeo/pzsvc-pdal/blob/master/LICENSE)

# pzsvc-pdal

Providing a [PDAL](http://pdal.io)-based microservice for Piazza.

The going in assumption is that we will receive some message from the [dispatcher](https://github.com/venicegeo/pz-dispatcher) indicating that a point cloud service has been requested. We will have the path to the data and a description of the task to be performed. We also have a responsibility to update the [job manager](https://github.com/venicegeo/pz-jobmanager) periodically with status updates.

## Install

For `pzsvc-pdal` to function properly, PDAL must be installed on your system. Our manifest.yml file specifies a custom buildpack to ensure that PDAL is available on Cloud Foundry. For local operation, follow installation instructions for your system, e.g., `brew install pdal` on Mac OS X

Go 1.5+ is required. You can download it [here](https://golang.org/dl/).

If you have not already done so, make sure you've setup your Go [workspace](https://golang.org/doc/code.html#Workspaces) and set the necessary environment [variables](https://golang.org/doc/code.html#GOPATH)

We make use of [Go 1.5's vendor/ experiment](https://medium.com/@freeformz/go-1-5-s-vendor-experiment-fd3e830f52c3#.ueuy8ao53), so you'll need to make sure you are running Go 1.5+, and that your `GO15VENDOREXPERIMENT` environment variable is set to `1`.

Installing `pzsvc-pdal` is as simple as

```console
$ export GO15VENDOREXPERIMENT=1
$ go get github.com/venicegeo/pzsvc-pdal
$ go install github.com/venicegeo/pzsvc-pdal
```

Assuming `$GOPATH/bin` is on your `$PATH`, the service can easily be started on port 8080

```console
$ pzsvc-pdal
```

## Examples

Perhaps the most straightforward means of demonstrating the `pzsvc-pdal` service is via [Postman](https://www.getpostman.com).

Begin by importing our [collection](https://github.com/venicegeo/pzsvc-pdal/blob/master/postman/pzsvc-pdal.json.postman_collection).

We also provide two environments, one to setup [localhost](https://github.com/venicegeo/pzsvc-pdal/blob/master/postman/pzsvc-pdal.json.postman_environment.local), another to setup [Cloud Foundry](https://github.com/venicegeo/pzsvc-pdal/blob/master/postman/pzsvc-pdal.json.postman_environment.cf).

## Testing

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

## Modifying

We use Godeps to aid in deployment. Upon saving, run

```console
$ godep save ./...
```

to update the Godeps folder and all import paths.

## Swagger

We have also begun to document the API via Swagger. The current API specification can be found [here](https://github.com/venicegeo/pzsvc-pdal/blob/master/swagger/swagger.yaml), but it is currently incomplete.

## Deploying

When deployed, `localhost:8080` is replaced with `pzsvc-pdal.cf.piazzageo.io`.

All commits to master will be pushed through the VeniceGeo DevOps infrastructure, first triggering a build in Jenkins and, upon success, pushing the resulting binaries to Cloud Foundry
