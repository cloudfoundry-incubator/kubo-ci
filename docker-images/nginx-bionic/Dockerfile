# Dockerfile for Ubuntu 18.04 nginx

# Pull base image.
FROM ubuntu:18.04

# Maintainer
MAINTAINER CFCR <cfcr@pivotal.io>

# Install Packages
RUN DEBIAN_FRONTEND=noninteractive apt-get update -y && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y \
  curl \
  nginx

RUN ln -sf /dev/stdout /var/log/nginx/access.log \
	&& ln -sf /dev/stderr /var/log/nginx/error.log

COPY ./nginx.conf /etc/nginx/sites-enabled/default

EXPOSE 443 80

STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]
