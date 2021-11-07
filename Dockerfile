FROM alpine:latest

MAINTAINER tzwsoho "tzwsoho@hotmail.com"

WORKDIR /

ADD TorrentCrawler /

ENV MYSQL_IP=127.0.0.1
ENV MYSQL_PORT=3306
ENV MYSQL_USERNAME=root
ENV MYSQL_PASSWORD=123456
ENV PATH=/:$PATH

CMD ["TorrentCrawler"]