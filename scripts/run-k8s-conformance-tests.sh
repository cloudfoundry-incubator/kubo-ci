#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

export GOPATH="${ROOT}/git-kubo-ci"
export CONFORMANCE_RELEASE_VERSION="$(cat kubo-version/version)"
export CONFORMANCE_RESULTS_DIR="${ROOT}/${CONFORMANCE_RESULTS_DIR}"

main() {
  source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
  setup_env "${KUBO_ENVIRONMENT_DIR}"

  local tmpfile="$(mktemp)" && echo "CONFIG=${tmpfile}"
  "${ROOT}/git-kubo-ci/scripts/generate-test-config.sh" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}" > "${tmpfile}"

  CONFIG="${tmpfile}" ginkgo -progress -v "${ROOT}/git-kubo-ci/src/tests/conformance"
}

main
