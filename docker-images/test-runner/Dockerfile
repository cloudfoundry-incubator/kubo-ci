FROM golang:1.9-alpine3.6
MAINTAINER pcf-kubo@pivotal.io

RUN apk --update  add git openssh bash && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

RUN go get github.com/onsi/ginkgo/ginkgo && go get github.com/onsi/gomega


