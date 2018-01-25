FROM alpine:latest
MAINTAINER pcf-kubo@pivotal.io

RUN apk add --update bash git jq
RUN wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.44-linux-amd64 -O bosh && \
  install bosh /usr/local/bin
