FROM alpine:3.6
MAINTAINER pcf-kubo@pivotal.io

RUN apk update && apk add \ 
      bash \
      curl \
      less \
      groff \
      jq \
      python \
      py-pip \
      py2-pip && \
      pip install --upgrade pip awscli && \
      mkdir /root/.aws

# BOSH CLI
RUN curl -O https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.42-linux-amd64 && \
  chmod +x bosh-cli-2.0.42-linux-amd64 && \
  mv bosh-cli-2.0.42-linux-amd64 /usr/local/bin/bosh-cli
