.PHONY: test get-deps

all: test

SKK_JISYO ?= http://openlab.jp/skk/dic/SKK-JISYO.L.gz

get-deps:
	go get -v -u ./...

test: SKK-JISYO.L
	go test -v ./...

SKK-JISYO.L:
	wget ${SKK_JISYO}
	gzip -d SKK-JISYO.L.gz
