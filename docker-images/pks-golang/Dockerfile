FROM golang:latest as builder

RUN apt-get update --yes && \
  apt-get install --yes git && \
  go get github.com/onsi/ginkgo/ginkgo

FROM golang:latest

COPY --from=builder /go/bin/ginkgo /go/bin/ginkgo
