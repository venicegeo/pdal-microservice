#!/bin/bash

APP=pzsvc-pdal
EXT=bin
SHA=$(git rev-parse HEAD)
SHORT=$(git rev-parse --short HEAD)
ARTIFACT="$SHA.$EXT"
S3BUCKET="venice-artifacts"
S3KEY=$APP/$ARTIFACT
S3URL=s3://$S3BUCKET/$S3KEY
STACK="$APP-$SHORT"
DOMAIN="cf.piazzageo.io"
