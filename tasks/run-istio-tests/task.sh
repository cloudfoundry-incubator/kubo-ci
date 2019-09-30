#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)

export KUBECONFIG=$ROOT/gcs-kubeconfig/config
ISTIO_VERSION=1.2.2

git clone --depth 50 --branch $ISTIO_VERSION https://github.com/istio/istio $GOPATH/src/istio.io/istio

cd $GOPATH/src/istio.io/istio

# init helm
./bin/init_helm.sh
kubectl create namespace istio-system

helm template install/kubernetes/helm/istio-init --name istio-init --namespace istio-system > istio-init.yml
trap "kubectl delete -f istio-init.yml; kubectl delete namespace istio-system" 0 1 2 3 15
kubectl apply -f istio-init.yml

timeout 15s ruby "$ROOT/git-kubo-ci/tasks/run-istio-tests/wait_for_apply_to_finish.rb" 23

helm template install/kubernetes/helm/istio --name istio-system --namespace istio-system --set global.mtls.enabled=true \
  --set sidecarInjectorWebhook.enabled=true --set global.hub=istio --set global.tag=$ISTIO_VERSION \
  --set global.crds=false > istio.yml
trap "kubectl delete -f istio.yml; kubectl delete -f istio-init.yml; kubectl delete namespace istio-system" 0 1 2 3 15
kubectl apply -f istio.yml

make e2e_simple TAG=${ISTIO_VERSION} E2E_ARGS='--installer=helm --skip_setup'
