FROM docker:dind
MAINTAINER pcf-kubo@pivotal.io

# # Packages
# RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y \
#   #gcc \
#   git-core \
#   wget

# # Ruby required for bosh-cli create-env
# RUN git clone https://github.com/postmodern/ruby-install.git /usr/local/ruby-install
# ENV PATH /usr/local/ruby-install/bin:/usr/local/bin:$PATH
# RUN ruby-install --system ruby 2.5.0

RUN apk update && \
    apk upgrade && \
    apk add curl wget bash git

# Install ruby and ruby-bundler
RUN apk add ruby ruby-bundler

# Clean APK cache
RUN rm -rf /var/cache/apk/*

# Direnv
RUN wget -O /usr/local/bin/direnv https://github.com/direnv/direnv/releases/download/v2.12.2/direnv.linux-amd64 && \
  chmod +x /usr/local/bin/direnv

# BOSH CLI
RUN wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.48-linux-amd64 -O bosh && \
  install bosh /usr/local/bin && \
  ln -s /usr/local/bin/bosh /usr/local/bin/bosh-cli

ADD common.sh /common.sh
