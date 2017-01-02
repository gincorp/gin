#!/usr/bin/env bash

set -e

if [ ! -z "${TRAVIS_TAG}" ] && [ "${TRAVIS_GO_VERSION}" = "1.7.4" ]; then
    go get -u github.com/goreleaser/releaser
    releaser

    cp -v dist/gin_Linux_x86_64/gin linux/gin

    docker login -u="${DOCKER_USERNAME}" -p="${DOCKER_PASSWORD}"
    docker build -t gincorp/gin:${TRAVIS_TAG} .
    docker build -t gincorp/gin .
    docker push gincorp/gin
fi
