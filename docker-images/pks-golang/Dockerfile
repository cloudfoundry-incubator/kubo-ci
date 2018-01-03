FROM golang:alpine as builder

RUN apk update && \
  apk add git && \
  go get github.com/onsi/ginkgo/ginkgo

FROM golang:alpine

RUN apk update && \
  apk add bash

COPY --from=builder /go/bin/ginkgo /go/bin/ginkgo
