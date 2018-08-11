.PHONY: test clean format deps build install all

all: clean deps build install

build:
	go build

install:
	go install ./...

test:
	go test -cover  ./...

format:
	gofmt -w ./

clean:
	go clean ./...

deps:
	go get google.golang.org/api/photoslibrary/v1
	go get golang.org/x/oauth2
	go get -u cloud.google.com/go/...
	go get github.com/nmrshll/oauth2-noserver
	go get -u github.com/utahta/go-openuri

