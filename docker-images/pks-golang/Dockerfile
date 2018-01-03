FROM golang:alpine as builder

RUN apk update && \
  apk add git && \
  go get github.com/onsi/ginkgo/ginkgo && \
  go get github.com/direnv/direnv

FROM golang:alpine

RUN apk update && \
  apk add bash

COPY --from=builder /go/bin/ginkgo /go/bin/direnv /go/bin/
