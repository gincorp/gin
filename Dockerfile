FROM alpine
MAINTAINER jspc <james@zero-internet.org.uk>

ADD linux/gin /
ENTRYPOINT ["/gin"]

EXPOSE 8000 8080
