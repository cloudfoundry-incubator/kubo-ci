#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
setup_env "${KUBO_ENVIRONMENT_DIR}"

go get istio.io/istio || true # go get returns error "no Go files", which is expected
cd $GOPATH/src/istio.io/istio
git checkout -b 0.8.0
make e2e_simple E2E_ARGS="--istioctl=/usr/local/istio-0.8.0/bin/istioctl" TAG='0.8.0'
