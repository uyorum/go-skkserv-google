FROM alpine:latest

ADD SKK-JISYO.L /
ADD skkserv-google /

EXPOSE 1178
CMD ["/skkserv-google", "-v", "-p", "1178", "/SKK-JISYO.L"]
