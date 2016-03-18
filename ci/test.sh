#!/bin/bash -ex

pushd `dirname $0`/.. > /dev/null
root=$(pwd -P)
popd > /dev/null

export GOPATH=$root/gogo
mkdir -p $GOPATH

source $root/ci/vars.sh

###
export GO15VENDOREXPERIMENT=1

go get -v github.com/venicegeo/$APP

go test -v $(go list github.com/venicegeo/$APP/... | grep -v /vendor/)

###
