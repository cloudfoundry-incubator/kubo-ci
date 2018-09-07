#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"

export GOPATH="${ROOT}/git-kubo-ci"

main() {
  local tmpfile="$(mktemp)" && echo "CONFIG=${tmpfile}"
  "${ROOT}/git-kubo-ci/scripts/generate-test-config.sh" > "${tmpfile}"

  CONFIG="${tmpfile}" ginkgo -progress -r "${ROOT}/git-kubo-ci/src/tests/turbulence-tests/"
}

main
