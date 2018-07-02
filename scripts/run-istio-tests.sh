#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
export PATH="$PATH:/root/go/out/linux_amd64/release"

source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
setup_env "${KUBO_ENVIRONMENT_DIR}"

go get istio.io/istio || true # go get returns error "no Go files", which is expected
cd $GOPATH/src/istio.io/istio
git checkout tags/0.8.0 -b 0.8.0
make e2e_simple TAG='0.8.0' E2E_ARGS='--installer=helm --skip_cleanup'
