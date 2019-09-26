#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

export KUBECONFIG=$ROOT/gcs-kubeconfig/config
ISTIO_VERSION=1.2.2

git clone --depth 50 --branch $ISTIO_VERSION https://github.com/istio/istio $GOPATH/src/istio.io/istio

cd $GOPATH/src/istio.io/istio

# init helm
./bin/init_helm.sh
kubectl apply -f install/kubernetes/helm/helm-service-account.yaml
helm init --service-account tiller
kubectl rollout status -w -n kube-system deployment.apps/tiller-deploy

# install Istio
trap "helm del --purge istio-init; kubectl delete namespace istio-system" 0 1 2 3 15
kubectl apply -f install/kubernetes/helm/helm-service-account.yaml
helm install install/kubernetes/helm/istio-init --name istio-init --namespace istio-system
trap "helm del --purge istio-system; kubectl delete namespace istio-system" 0 1 2 3 15
helm install install/kubernetes/helm/istio --name istio-system --namespace istio-system --set global.mtls.enabled=true \
  --set sidecarInjectorWebhook.enabled=true --set global.hub=istio --set global.tag=$ISTIO_VERSION --set global.crds=false --debug

make e2e_simple TAG=${ISTIO_VERSION} E2E_ARGS='--installer=helm --skip_setup'
