.PHONY: test prepare

SKK_JISYO ?= http://openlab.jp/skk/dic/SKK-JISYO.L.gz

test: SKK-JISYO.L
	go test .

SKK-JISYO.L:
	wget ${SKK_JISYO}
	gzip -d SKK-JISYO.L.gz
