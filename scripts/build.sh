#!/bin/bash

echo "Compiling for linux..."
GOOS=linux GOARCH=amd64 go build .

echo "Constructing Docker image"
docker build -t chambbj/pzsvc-pdal .
docker push chambbj/pzsvc-pdal

echo "Cleaning up..."
rm pzsvc-pdal
