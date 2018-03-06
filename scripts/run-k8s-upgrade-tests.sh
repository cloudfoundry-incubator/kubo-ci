#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

export GOPATH="${ROOT}/git-kubo-ci"
export NEW_BOSH_STEMCELL_VERSION="$(cat ${ROOT}/new-bosh-stemcell/version)"

setup_env() {
  mkdir -p "${KUBO_ENVIRONMENT_DIR}"
  cp "${ROOT}/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "${ROOT}/kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

  "${ROOT}/git-kubo-deployment/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"
  "${ROOT}/git-kubo-deployment/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}"
}

main() {
  setup_env

  local tmpfile="$(mktemp)" && echo "CONFIG=${tmpfile}"
  "${ROOT}/git-kubo-ci/scripts/generate-test-config.sh" ${KUBO_ENVIRONMENT_DIR} ${DEPLOYMENT_NAME} > "${tmpfile}"

  CONFIG="${tmpfile}" ginkgo -r -progress -v "${ROOT}/git-kubo-ci/src/tests/upgrade-tests/"
}

main
