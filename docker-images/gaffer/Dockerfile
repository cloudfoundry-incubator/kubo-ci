FROM node:8.9-alpine

RUN apk add --update python python-dev make g++ wget unzip ca-certificates && rm -rf /var/cache/apk/*
RUN wget https://storage.googleapis.com/kubo-public/services-gaffer.zip?ignoreCache=1 -O services-gaffer.zip && \
  unzip services-gaffer.zip && \
  cd services-gaffer-master && \
  npm install --only production