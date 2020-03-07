  
export GO111MODULE=on
VERSION=$(shell git describe --tags --candidates=1 --dirty)
IMG ?= jdamata/k8s-events:${VERSION}

.PHONY: run

run:
	go run main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

docker-build: test
	docker build . -t ${IMG}

docker-push:
	docker push ${IMG}
