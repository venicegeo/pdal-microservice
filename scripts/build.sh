#!/bin/bash

echo "Compiling for linux..."
GOOS=linux GOARCH=amd64 go build .

echo "Constructing Docker image"
docker build -t chambbj/cflinuxfs2-pdal .
docker push chambbj/cflinuxfs2-pdal

echo "Cleaning up..."
rm pzsvc-pdal
