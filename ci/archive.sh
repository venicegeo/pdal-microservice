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

go install -v github.com/venicegeo/$APP


###

src=$GOPATH/bin/$APP

# stage the artifact for a mvn deploy
mv $src $root/$APP.$EXT
