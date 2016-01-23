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
Package s3 provides a number of S3 helper functions.

For example,

	// Get the filename from the key.
	filename := s3.ParseFilenameFromKey(key)

	// Create the file.
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Download the data.
	err = s3.Download(file, bucket, key)
	if err != nil {
		panic(err)
	}
*/
package s3

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

/*
Download downloads a file from S3.

This is merely a wrapper around the aws-sdk-go downloader. It allows us to
isolate the aws-sdk-go dependencies and unify error handling.
*/
func Download(file *os.File, bucket, key string) error {
	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	log.Println("Downloaded", numBytes, "bytes")
	return nil
}

/*
Upload uploads a file to S3.

This is merely a wrapper around the aws-sdk-go uploader. It allows us to isolate
the aws-sdk-go dependencies and unify error handling.
*/
func Upload(file *os.File, bucket, key string) error {
	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	log.Println("Successfully uploaded to", result.Location)
	return nil
}

/*
ParseFilenameFromKey parses the S3 filename from the key.

This would typically be used prior to s3.Download or s3.Upload to create the
required *os.File.
*/
func ParseFilenameFromKey(key string) string {
	keySlice := strings.Split(key, "/")
	return keySlice[len(keySlice)-1]
}
