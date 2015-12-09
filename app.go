package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Hello)
	http.Handle("/", r)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Hello(w http.ResponseWriter, req *http.Request) {
	file, err := os.Create("download_file.laz")
	if err != nil {
		fmt.Fprintln(w, "Failed to create file", err)
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String("venicegeo-sample-data"),
			Key:    aws.String("pointcloud/samp11-utm.laz"),
		})
	if err != nil {
		fmt.Fprintln(w, "Failed to download file", err)
		return
	}

	fmt.Fprintln(w, "Downloaded file", file.Name(), numBytes, "bytes")

	out, _ := exec.Command("pdal", "info", file.Name()).CombinedOutput()
	if err != nil {
		fmt.Fprintln(w, err)
	}
	fmt.Fprintln(w, string(out))
}
