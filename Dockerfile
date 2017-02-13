FROM alpine:latest

ADD SKK-JISYO.L /
ADD skkserv-google /

RUN mkdir /lib64 && ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

EXPOSE 1178
CMD ["/skkserv-google", "-v", "-p", "1178", "/SKK-JISYO.L"]
