#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

kubeconfig="gcs-kubeconfig/${KUBECONFIG_FILE}"
export GOPATH="${ROOT}/git-kubo-ci"

main() {
  if [[ ! -e "${kubeconfig}" ]]; then
    echo "Did not find kubeconfig at gcs-kubeconfig/${KUBECONFIG_FILE}!"
    exit 1
  fi

  mkdir -p ~/.kube
  cp ${kubeconfig} ~/.kube/config

  source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
  create_environment_dir "${KUBO_ENVIRONMENT_DIR}"

  local tmpfile="$(mktemp)" && echo "CONFIG=${tmpfile}"
  "${ROOT}/git-kubo-ci/scripts/generate-test-config.sh" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}" > "${tmpfile}"

  CONFIG="${tmpfile}" ginkgo -failFast -progress -r "${ROOT}/git-kubo-ci/src/tests/turbulence-tests/"
}

main
