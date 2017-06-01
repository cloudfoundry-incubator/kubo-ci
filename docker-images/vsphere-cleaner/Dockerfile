FROM bash:latest
MAINTAINER pcf-kubo@pivotal.io

RUN apk update && apk add ca-certificates && update-ca-certificates && apk add openssl && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

# vsphere-cleaner
RUN wget https://storage.googleapis.com/kubo-public/vsphere-cleaner -O /usr/bin/vsphere-cleaner && \
  chmod +x /usr/bin/vsphere-cleaner
