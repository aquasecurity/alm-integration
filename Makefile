# Set the default goal
.DEFAULT_GOAL := build

# Active module mode, as we use Go modules to manage dependencies
export GO111MODULE=on

GO_FMT=gofmt

.PHONY: all build fmt vet test

default : all

.PHONY: build
build :
	@echo "building...."
	CGO_ENABLED=0 go build -o ./bin/postee main.go
	@echo "Done!"

docker :
	@echo "Building image...."
	docker build -t aquasec/postee:latest -f Dockerfile .

fmt :
	@echo "fmt...."
	$(GO_FMT) -s -w ./

test :
	go test -race -coverprofile=coverage.txt -covermode=atomic ./alertmgr ./scanservice ./dbservice ./formatting ./data ./regoservice

cover :
	go test ./scanservice ./dbservice ./alertmgr ./formatting ./data ./regoservice -v -coverprofile=cover.out
	go tool cover -html=cover.out
