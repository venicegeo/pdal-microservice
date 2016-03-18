#!/bin/bash -ex

pushd `dirname $0`/.. > /dev/null
root=$(pwd -P)
popd > /dev/null

# Run the test
newman -c $root/postman/pzsvc-pdal-black-box-tests.json.postman_collection -e $root/postman/pzsvc-pdal.json.postman_environment.cf -s
