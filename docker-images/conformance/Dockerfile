FROM golang:1.10.3-alpine3.8
MAINTAINER pcf-kubo@pivotal.io

RUN apk add --update \
  git \
  bash \
  curl

ENV K8S_VERSION="v1.10.5"

RUN go get -d k8s.io/kubernetes/test/... && \
  go get k8s.io/test-infra/kubetest

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl && \
  chmod +x ./kubectl && \
  mv ./kubectl /usr/local/bin/kubectl

ENV KUBECTL_PATH=/usr/local/bin/kubectl
ENV KUBERNETES_CONFORMANCE_TEST=y

WORKDIR $GOPATH/src/k8s.io/kubernetes
RUN git checkout tags/${K8S_VERSION}
