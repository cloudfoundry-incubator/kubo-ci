FROM ubuntu:14.04
MAINTAINER pcf-kubo@pivotal.io

# Packages
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y \
  curl \
  git-core \
  jq

WORKDIR /tmp/docker-build


# Ruby required for bosh-cli create-env
RUN git clone https://github.com/postmodern/ruby-install.git /usr/local/ruby-install
ENV PATH /usr/local/ruby-install/bin:/usr/local/bin:$PATH
RUN ruby-install --system ruby 2.4.1

COPY Gemfile .
COPY Gemfile.lock .
RUN gem install bundler --no-ri --no-rdoc
RUN bundle install

