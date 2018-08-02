#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
setup_env "${KUBO_ENVIRONMENT_DIR}"
ISTIO_VERSION=1.0.0

git clone --depth 50 --branch $ISTIO_VERSION https://github.com/istio/istio $GOPATH/src/istio.io/istio

cd $GOPATH/src/istio.io/istio
./bin/init_helm.sh
make e2e_simple TAG=${ISTIO_VERSION} E2E_ARGS='--installer=helm --use-automatic-injection'
