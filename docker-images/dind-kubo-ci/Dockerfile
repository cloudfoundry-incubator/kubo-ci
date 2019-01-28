FROM docker:dind
MAINTAINER pcf-kubo@pivotal.io

RUN apk add --no-cache curl wget bash git ruby ruby-bundler

# BOSH CLI
RUN wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-5.4.0-linux-amd64 -O bosh && \
  install bosh /usr/local/bin && \
  ln -s /usr/local/bin/bosh /usr/local/bin/bosh-cli
