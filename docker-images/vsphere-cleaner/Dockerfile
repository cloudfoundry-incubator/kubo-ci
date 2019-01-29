FROM bash:latest
MAINTAINER pcf-kubo@pivotal.io

RUN apk add --no-cache ca-certificates && update-ca-certificates && apk add --no-cache openssl

# vsphere-cleaner
RUN wget https://storage.googleapis.com/kubo-public/vsphere-cleaner -O /usr/bin/vsphere-cleaner && \
  chmod +x /usr/bin/vsphere-cleaner
