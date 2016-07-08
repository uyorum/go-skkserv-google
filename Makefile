.PHONY: test get-deps

## Prepare and Run tests
all: test

SKK_JISYO ?= http://openlab.jp/skk/dic/SKK-JISYO.L.gz

## Install dependencies
get-deps:
	go get -v -u ./...

## Run tests
test: SKK-JISYO.L
	go test -v ./...

## Download SKK-JISYO.L
SKK-JISYO.L:
	wget ${SKK_JISYO}
	gzip -d SKK-JISYO.L.gz
