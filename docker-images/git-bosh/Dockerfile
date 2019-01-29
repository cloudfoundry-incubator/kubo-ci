FROM alpine:latest
MAINTAINER pcf-kubo@pivotal.io

RUN apk add --no-cache bash git jq
RUN wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-5.4.0-linux-amd64 -O bosh && \
  install bosh /usr/local/bin
