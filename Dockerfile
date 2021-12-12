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

ENV DEBIAN_FRONTEND=noninteractive
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

COPY --from=build /GoTorrentCrawler/TorrentCrawler /GoTorrentCrawler/TorrentCrawler

CMD ["TorrentCrawler"]
