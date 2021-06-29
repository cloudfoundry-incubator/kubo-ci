FROM ubuntu:16.04
MAINTAINER pcf-kubo@pivotal.io

# Packages
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y \
  bc \
  curl \
  gcc \
  jq \
  git-core \
  netcat-openbsd  \
  libcurl3  \
  make \
  python3-pip \
  python3-dev \
  python3-software-properties \
  software-properties-common \
  wget \
  ipcalc \
  zip \
  vim \
  haproxy \
  libssl-dev \
  libssl-doc \
  iptables # required for sshuttle
  # libssl required for Ruby

WORKDIR /tmp/docker-build

# Golang
ENV GOLANG_VERSION=1.13.3
RUN wget -q https://storage.googleapis.com/golang/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && rm go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV GOPATH /root/go
RUN mkdir -p /root/go/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin
RUN go get github.com/onsi/gomega && \
   go get github.com/onsi/ginkgo/ginkgo && \
   go get github.com/tsenart/vegeta

# Google SDK
ENV GCLOUD_VERSION=224.0.0
ENV GCLOUD_SHA1SUM=0a85de5c35c562f5d602ad20f567d45a214e91e5570ae95a93ced8361fa6d021

RUN wget -q https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GCLOUD_VERSION}-linux-x86_64.tar.gz \
    -O gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz && \
    echo "${GCLOUD_SHA1SUM}  gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz" > gcloud_${GCLOUD_VERSION}_SHA1SUM && \
    shasum -a 256 -cw --status gcloud_${GCLOUD_VERSION}_SHA1SUM && \
    tar xvf gcloud_${GCLOUD_VERSION}_linux_amd64.tar.gz && \
    mv google-cloud-sdk / && cd /google-cloud-sdk  && ./install.sh

ENV PATH=$PATH:/google-cloud-sdk/bin

# Ruby required for bosh-cli create-env
RUN git clone https://github.com/postmodern/ruby-install.git /usr/local/ruby-install
RUN /usr/local/ruby-install/bin/ruby-install --system ruby 2.5.3

# Credhub
RUN wget -q https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/2.2.0/credhub-linux-2.2.0.tgz \
  -O credhub-linux-2.2.0.tgz && tar xvf credhub-linux-2.2.0.tgz && mv credhub /usr/bin

# kubectl
ENV KUBE_VERSION="v1.18.20"
ADD https://storage.googleapis.com/kubernetes-release/release/$KUBE_VERSION/bin/linux/amd64/kubectl /usr/bin/kubectl
RUN chmod +x /usr/bin/kubectl

# BOSH CLI
RUN wget -q https://github.com/cloudfoundry/bosh-cli/releases/download/v6.3.1/bosh-cli-6.3.1-linux-amd64 -O bosh && \
  echo "641f84b14b1c1d6a39a949526c5b257821090effbc5e8d8bb8c528660258cdf9  bosh" > bosh-shasum && \
  shasum -a 256 -cw --status bosh-shasum && \
  install bosh /usr/local/bin

# Openstack CLI
RUN pip3 install -U setuptools
RUN pip3 install cryptography==2.9.2
RUN pip3 install pyOpenSSL==19.1.0
RUN pip3 install python-glanceclient==3.2.0
RUN pip3 install python-openstackclient==5.3.0

# AWS CLI
RUN pip3 install awscli

# sshuttle requires trunk due to a bug with Python 3 recently fixed
#RUN pip3 install sshuttle
RUN git clone https://github.com/sshuttle/sshuttle.git && cd sshuttle && python3 setup.py install

# GOVC CLI
RUN wget -q -O - -o /dev/null https://github.com/vmware/govmomi/releases/download/v0.17.1/govc_linux_amd64.gz | gunzip > /usr/local/bin/govc && \
  chmod +x /usr/local/bin/govc

RUN gem install bundler --no-ri --no-rdoc && gem update --system

ARG SPRUCE_VERSION=v1.16.2
RUN wget -q https://github.com/geofffranks/spruce/releases/download/${SPRUCE_VERSION}/spruce-linux-amd64 -O /usr/bin/spruce && \
  chmod +x /usr/bin/spruce

RUN git clone https://github.com/fsaintjacques/semver-tool && \
  cd semver-tool && make install

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
