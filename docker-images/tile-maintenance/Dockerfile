FROM ubuntu:14.04
MAINTAINER pcf-kubo@pivotal.io

# Packages
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y \
  curl \
  git-core \
  wget \
  jq \
  python-pip \
  python-dev \
&& rm -rf /var/lib/apt/lists/*

WORKDIR /tmp/docker-build

# Google SDK
ENV GCLOUD_VERSION=144.0.0
ENV GCLOUD_SHA1SUM=4ef237d6d2696b353a5a1ef927462ef3833d261f

RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_VERSION}-linux-x86_64.tar.gz \
    -O gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz && \
    echo "${GCLOUD_SHA1SUM}  gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz" > gcloud_${GCLOUD_VERSION}_SHA1SUM && \
    sha1sum -cw --status gcloud_${GCLOUD_VERSION}_SHA1SUM && \
    tar xvf gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz && \
    mv google-cloud-sdk / && cd /google-cloud-sdk  && ./install.sh

ENV PATH=$PATH:/google-cloud-sdk/bin

# Ruby required for bosh-cli create-env
RUN git clone https://github.com/postmodern/ruby-install.git /usr/local/ruby-install
ENV PATH /usr/local/ruby-install/bin:/usr/local/bin:$PATH
RUN ruby-install --system ruby 2.4.1

# BOSH CLI
RUN wget https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.44-linux-amd64 && \
  chmod +x bosh-cli-2.0.44-linux-amd64 && \
  mv bosh-cli-2.0.44-linux-amd64 /usr/local/bin/bosh-cli

COPY Gemfile .
COPY Gemfile.lock .
RUN gem install bundler --no-ri --no-rdoc
RUN bundle install

