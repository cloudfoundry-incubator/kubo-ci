#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
setup_env "${KUBO_ENVIRONMENT_DIR}"
ISTIO_VERSION=1.0.0

git clone --depth 50 --branch $ISTIO_VERSION https://github.com/istio/istio $GOPATH/src/istio.io/istio

cd $GOPATH/src/istio.io/istio

# init helm
./bin/init_helm.sh
kubectl apply -f install/kubernetes/helm/helm-service-account.yaml
helm init --service-account tiller
kubectl rollout status -w -n kube-system deployment.apps/tiller-deploy

# install Istio
kubectl apply -f install/kubernetes/helm/istio/templates/crds.yaml
kubectl apply -f install/kubernetes/helm/istio/charts/certmanager/templates/crds.yaml
trap "helm del --purge istio-system; kubectl delete namespace istio-system" 0 1 2 3 15
helm install install/kubernetes/helm/istio --name istio-system --namespace istio-system --set global.mtls.enabled=true \
  --set sidecarInjectorWebhook.enabled=true --set global.hub=istio --set global.tag=$ISTIO_VERSION --set global.crds=false --debug

make e2e_simple TAG=${ISTIO_VERSION} E2E_ARGS='--installer=helm --skip_setup'
