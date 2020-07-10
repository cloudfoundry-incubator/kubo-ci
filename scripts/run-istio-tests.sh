#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

export KUBECONFIG=$ROOT/gcs-kubeconfig/config
ISTIO_VERSION=1.5.8

git clone --depth 50 --branch $ISTIO_VERSION https://github.com/istio/istio $GOPATH/src/istio.io/istio

cd $GOPATH/src/istio.io/istio

# init helm - removed upstream , pasting here for now prior to migrating to istioctl

export GO_TOP=${GO_TOP:-$(echo "${GOPATH}" | cut -d ':' -f1)}
export OUT_DIR=${OUT_DIR:-${GO_TOP}/out}

HELM_VER=${HELM_VER:-v2.10.0}

export GOPATH=${GOPATH:-$GO_TOP}
# Normally set by Makefile
export ISTIO_BIN=${ISTIO_BIN:-${GOPATH}/bin}

# Set the architecture. Matches logic in the Makefile.
export GOARCH=${GOARCH:-'amd64'}

# Determine the OS. Matches logic in the Makefile.
LOCAL_OS=${OSTYPE}
case $LOCAL_OS in
  "linux"*)
    LOCAL_OS='linux'
    ;;
  "darwin"*)
    LOCAL_OS='darwin'
    ;;
  *)
    echo "This system's OS ${LOCAL_OS} isn't recognized/supported"
    exit 1
    ;;
esac
export GOOS=${GOOS:-${LOCAL_OS}}

# Gets the download command supported by the system (currently either curl or wget)
# simplify because init.sh has already decided which one is
DOWNLOAD_COMMAND=""
function set_download_command () {
    # Try curl.
    if command -v curl > /dev/null; then
        if curl --version | grep Protocols  | grep https > /dev/null; then
	       DOWNLOAD_COMMAND='curl -Lo '
	       return
        fi
    fi

    # Try wget.
    if command -v wget > /dev/null; then
        DOWNLOAD_COMMAND='wget -O '
        return
    fi
}
set_download_command

# test scripts seem to like to run this script directly rather than use make
export ISTIO_OUT=${ISTIO_OUT:-${ISTIO_BIN}}

# install helm if not present, it must be the local version.
if [ ! -f "${ISTIO_OUT}/version.helm.${HELM_VER}" ] ; then
    TD=$(mktemp -d)
    # Install helm. Please keep it in sync with .circleci
    cd "${TD}" && \
        ${DOWNLOAD_COMMAND} "${TD}/helm.tgz" "https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VER}-${LOCAL_OS}-amd64.tar.gz" && \
        tar xfz helm.tgz && \
        mv ${LOCAL_OS}-amd64/helm "${ISTIO_OUT}/helm-${HELM_VER}" && \
        cp "${ISTIO_OUT}/helm-${HELM_VER}" "${ISTIO_OUT}/helm" && \
        rm -rf "${TD}" && \
        touch "${ISTIO_OUT}/version.helm.${HELM_VER}"
fi

cd $GOPATH/src/istio.io/istio

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
