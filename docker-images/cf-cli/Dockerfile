FROM alpine:latest
RUN apk add --update openssl && \
  rm -rf /var/cache/apk/* && \
  wget 'https://cli.run.pivotal.io/stable?release=linux64-binary&version=6.30.0&source=github-rel' -O cf.tgz && \
  tar -xvf cf.tgz && \
  rm cf.tgz && \
  chmod a+x cf && \
  mv cf /usr/bin
