FROM golang:1.11-alpine3.8
MAINTAINER pcf-kubo@pivotal.io

RUN apk --no-cache add git openssh bash
ENV CGO_ENABLED=0

RUN go get github.com/onsi/ginkgo/ginkgo && go get github.com/onsi/gomega
