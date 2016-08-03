# go-skkserv-google
A SKK server using Google conversion API.

[![Build Status](https://travis-ci.org/uyorum/go-skkserv-google.svg?branch=master)](https://travis-ci.org/uyorum/go-skkserv-google)

This uses Google API only when "okuri-nasi" and when "okuri-ari" this usues SKK dictionary file.

## How to use

``` shell
$ make SKK-JISYO.L
$ make get-deps
$ make build
$ ./skkserv-google ./SKK-JISYO.L
```

### Options

``` shell
-p int
      Port number skkserv uses (default 1178)
```

### Run in Docker

``` shell
$ docker run -d -p 127.0.0.1:1178:1178 uyorum/skkserv-google
```
