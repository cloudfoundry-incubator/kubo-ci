FROM alpine:latest
RUN apk add --update curl && rm -rf /var/cache/apk/* && curl  -LO https://storage.googleapis.com/kubernetes-release/release/v1.7.1/bin/linux/amd64/kubectl && chmod a+x kubectl && mv kubectl /usr/bin
