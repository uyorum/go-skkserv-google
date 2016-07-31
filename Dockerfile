FROM alpine:latest

MAINTAINER uyorum uyorum.pub@gmail.com

ADD SKK-JISYO.L /
ADD skkserv-google /

EXPOSE 1178
CMD ["/skkserv-google", "-p", "1178", "/SKK-JISYO.L"]
