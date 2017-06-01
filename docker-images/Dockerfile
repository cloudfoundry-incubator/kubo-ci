FROM ubuntu:14.04
MAINTAINER pcf-kubo@pivotal.io

# Packages
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y \
  gcc \
  git-core \
  make \
  python-software-properties \
  software-properties-common \
  wget

WORKDIR /tmp/docker-build

# Direnv
RUN wget https://github.com/direnv/direnv/releases/download/v2.9.0/direnv.linux-amd64 -O /usr/local/bin/direnv
RUN chmod +x /usr/local/bin/direnv
RUN echo 'eval "$(direnv hook bash)"' >> ~/.bashrc

# Golang
RUN wget https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.7.4.linux-amd64.tar.gz
ENV GOPATH /root/go
RUN mkdir -p /root/go/bin
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin
RUN go get github.com/onsi/gomega
RUN go get github.com/onsi/ginkgo
RUN go get github.com/cloudfoundry/bosh-cli
RUN go install github.com/onsi/ginkgo/...

# CF
RUN wget -O /cf-cli_amd64.deb https://cli.run.pivotal.io/stable?release=debian64
RUN dpkg -i /cf-cli_amd64.deb && rm /cf-cli_amd64.deb

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
RUN ruby-install --system ruby 2.3.3

# Credhub
RUN wget https://github.com/pivotal-cf/credhub-cli/releases/download/0.4.0/credhub-linux-0.4.0.tgz \
  -O credhub-linux-0.4.0.tgz && tar xvf credhub-linux-0.4.0.tgz && mv credhub /usr/bin

RUN wget https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/0.6.0/credhub-linux-0.6.0.tgz \
  -O credhub-linux-0.6.0.tgz && tar xvf credhub-linux-0.6.0.tgz && mv credhub /usr/bin/credhub-0.6.0

# kubectl
RUN wget https://storage.googleapis.com/kubernetes-release/release/v1.6.1/bin/linux/amd64/kubectl -O /usr/bin/kubectl && \
  chmod +x /usr/bin/kubectl

# Cleanup
RUN rm -rf /tmp/docker-build

# vsphere-cleaner
RUN wget https://storage.googleapis.com/kubo-public/vsphere-cleaner -O /usr/bin/vsphere-cleaner && \
  chmod +x /usr/bin/vsphere-cleaner
