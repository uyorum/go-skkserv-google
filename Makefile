.PHONY: test get-deps build docker

SKK_JISYO ?= http://openlab.jp/skk/dic/SKK-JISYO.L.gz
DOCKER_IMAGE ?= uyorum/skkserv-google

## Run tests
test:
	go test -v ./...

## Build a binary
build:
	go build skkserv-google.go

## Install dependencies
get-deps: SKK-JISYO.L
	go get -v -u ./...

## Download SKK-JISYO.L
SKK-JISYO.L:
	wget ${SKK_JISYO}
	gzip -d SKK-JISYO.L.gz

## Build docker image
docker: SKK-JISYO.L
	GOOS=linux GOARCH=amd64 go build skkserv-google.go
	docker build -t ${DOCKER_IMAGE} .
	rm -f skkserv-google
