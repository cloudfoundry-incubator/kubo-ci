FROM golang:1 as builder
MAINTAINER cfcr@pivotal.io

RUN apt-get update \
			&& apt-get install -y --no-install-recommends \
			rsync \
			&& rm -rf /var/lib/apt/lists/*

ENV KUBE_VERSION="v1.14.1"

RUN git clone --depth 50 --branch $KUBE_VERSION \
    https://github.com/kubernetes/kubernetes.git /go/src/k8s.io/kubernetes

ADD https://storage.googleapis.com/kubernetes-release/release/$KUBE_VERSION/bin/linux/amd64/kubectl /usr/bin/kubectl
RUN chmod +x /usr/bin/kubectl

ENV KUBECTL_PATH=/usr/local/bin/kubectl
ENV KUBERNETES_CONFORMANCE_TEST=y

WORKDIR $GOPATH/src/k8s.io/kubernetes

RUN make ginkgo
RUN make WHAT='test/e2e/e2e.test'

FROM ubuntu:16.04
MAINTAINER cfcr@pivotal.io

RUN apt-get update \
&& apt-get install -y --no-install-recommends \
ca-certificates \
vim \
&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/k8s.io/kubernetes/_output/bin/e2e.test /e2e.test
COPY --from=builder /go/src/k8s.io/kubernetes/_output/bin/ginkgo /usr/bin/ginkgo
COPY --from=builder /usr/bin/kubectl /usr/bin/kubectl

ENV KUBECTL_PATH /usr/bin/kubectl
ENV KUBECONFIG /kubeconfig

