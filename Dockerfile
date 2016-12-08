FROM alpine
MAINTAINER jspc <james@zero-internet.org.uk>

ADD linux/workflow-engine /
ENTRYPOINT ["/workflow-engine"]
