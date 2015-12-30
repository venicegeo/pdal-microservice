docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.4.2 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o pdal'
docker build -t venicegeo/pzsvc-pdal .
docker run --rm -it -p 8080:8080 -v ~/.aws/credentials:/root/.aws/credentials venicegeo/pzsvc-pdal
