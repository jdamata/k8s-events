  
export GO111MODULE=on
VERSION=$(shell git describe --tags --candidates=1 --dirty)
BUILD_FLAGS=-ldflags="-X main.version=$(VERSION)"
# CERT_ID ?= TODO
SRC=$(shell find . -name '*.go')

.PHONY: all clean release install

all: k8s-events-linux-amd64 k8s-events-darwin-amd64

clean:
	rm -f k8s-events k8s-events-linux-amd64 k8s-events-darwin-amd64

k8s-events-linux-amd64: $(SRC)
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ .

k8s-events-darwin-amd64: $(SRC)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ .

install:
	rm -f k8s-events
	go build $(BUILD_FLAGS) .
	mv k8s-events ~/bin/