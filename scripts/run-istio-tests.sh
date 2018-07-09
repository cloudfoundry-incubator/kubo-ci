#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
setup_env "${KUBO_ENVIRONMENT_DIR}"

go get -d github.com/cfcr/istio
mkdir -p $GOPATH/src/istio.io
ln -s $GOPATH/src/github.com/cfcr/istio $GOPATH/src/istio.io/istio
cd $GOPATH/src/istio.io/istio
git checkout 0.8.0
./bin/init_helm.sh
make e2e_simple TAG='0.8.0' E2E_ARGS='--installer=helm --use_automatic_injection'
