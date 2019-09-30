#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

export KUBECONFIG=$ROOT/gcs-kubeconfig/config
ISTIO_VERSION=1.2.2

git clone --depth 50 --branch $ISTIO_VERSION https://github.com/istio/istio $GOPATH/src/istio.io/istio

cd $GOPATH/src/istio.io/istio

# init helm
./bin/init_helm.sh
kubectl create namespace istio-system

helm template install/kubernetes/helm/istio-init --name istio-init --namespace istio-system > istio-init.yml
trap "kubectl delete -f istio-init.yml --ignore-not-found=true; kubectl delete namespace istio-system" 0 1 2 3 15
kubectl apply -f istio-init.yml

kubectl wait --for condition=complete --timeout=60s --all job -n istio-system
kubectl wait --for condition=established --timeout=60s --all crd

helm template install/kubernetes/helm/istio --name istio-system --namespace istio-system --set global.mtls.enabled=true \
  --set sidecarInjectorWebhook.enabled=true --set global.hub=istio --set global.tag=$ISTIO_VERSION \
  --set global.crds=false > istio.yml
trap "kubectl delete -f istio.yml --ignore-not-found=true; kubectl delete -f istio-init.yml --ignore-not-found=true; kubectl delete namespace istio-system" 0 1 2 3 15
kubectl apply -f istio.yml

make e2e_simple TAG=${ISTIO_VERSION} E2E_ARGS='--installer=helm --skip_setup'
