FROM alpine:latest
RUN apk add --update openssl jq && \
  rm -rf /var/cache/apk/* && \
  wget 'https://cli.run.pivotal.io/stable?release=linux64-binary&version=6.30.0&source=github-rel' -O cf.tgz && \
  tar -xvf cf.tgz && \
  rm cf.tgz && \
  chmod a+x cf && \
  mv cf /usr/bin && \
  wget 'https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.31-linux-amd64' -O bosh-cli && \
  chmod a+x bosh-cli && \
  mv bosh-cli /usr/bin
