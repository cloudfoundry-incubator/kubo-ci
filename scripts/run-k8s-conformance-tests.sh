#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail


generate_config() {
  ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
  KUBO_ENVIRONMENT_DIR="${ROOT}/environment"
  DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"

  export GOPATH="${ROOT}/git-kubo-ci"
  export CONFORMANCE_RELEASE_VERSION="$(cat kubo-version/version)"
  export CONFORMANCE_RESULTS_DIR="${ROOT}/${CONFORMANCE_RESULTS_DIR}"
  source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
  setup_env "${KUBO_ENVIRONMENT_DIR}"

  local tmpfile="$(mktemp)" && echo "CONFIG=${tmpfile}"
  export CONFIG="${tmpfile}"
  "${ROOT}/git-kubo-ci/scripts/generate-test-config.sh" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}" > "${CONFIG}"
}

main() {
  if [[ -z "${CONFIG+x}" ]]; then
    generate_config
  fi

  ginkgo -progress -v "${GOPATH}/src/tests/conformance"
}

main
