.PHONY: test get-deps

SKK_JISYO ?= http://openlab.jp/skk/dic/SKK-JISYO.L.gz

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
