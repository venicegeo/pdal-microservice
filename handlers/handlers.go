package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/venicegeo/pdal-microservice/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pdal-microservice/objects"
)

// var validPath = regexp.MustCompile("^/(info|pipeline)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(pdal)$")

// PdalHandler handles PDAL jobs.
func PdalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check that we have a valid path. Is this the correct place to do this?
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	// Parse the incoming JSON body, and unmarshal as events.NewData struct.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var msg objects.JobInput
	if err := json.Unmarshal(b, &msg); err != nil {
		log.Fatal(err)
	}

	var res objects.JobOutput
	res.Input = msg
	res.StartedAt = time.Now()
	res.Status = "started"

	file, err := os.Create("download_file.laz")
	if err != nil {
		// errors here should also be JSON-encoded as below
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(msg.Source.Bucket),
			Key:    aws.String(msg.Source.Key),
		})
	if err != nil {
		// errors here should also be JSON-encoded as below
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Code(), awsErr.Message())
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	log.Println("Downloaded", numBytes, "bytes")

	out, _ := exec.Command("pdal", msg.Function, file.Name()).CombinedOutput()

	// Trim whitespace
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, out); err != nil {
		fmt.Println(err)
	}

	if err = json.Unmarshal(buffer.Bytes(), &res.Response); err != nil {
		log.Fatal(err)
	}
	res.Status = "finished"
	res.FinishedAt = time.Now()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
}
