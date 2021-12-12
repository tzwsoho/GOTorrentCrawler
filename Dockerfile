FROM golang as build

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /GoTorrentCrawler
COPY . .

RUN go build -tags netgo

FROM alpine

MAINTAINER tzwsoho "tzwsoho@hotmail.com"

WORKDIR /GoTorrentCrawler
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY --from=build /GoTorrentCrawler/TorrentCrawler /GoTorrentCrawler/TorrentCrawler
COPY --from=build /GoTorrentCrawler/config.json /GoTorrentCrawler/config.json

CMD ["./TorrentCrawler"]
