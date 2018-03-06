#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

export GOPATH="${ROOT}/git-kubo-ci"
export ENABLE_ADDONS_TESTS="${ENABLE_ADDONS_TESTS}"
export ENABLE_API_EXTENSIONS_TESTS="${ENABLE_API_EXTENSIONS_TESTS}"
export ENABLE_GENERIC_TESTS="${ENABLE_GENERIC_TESTS}"
export ENABLE_IAAS_K8S_LB_TESTS="${ENABLE_IAAS_K8S_LB_TESTS}"
export ENABLE_MULTI_AZ_TESTS="${ENABLE_MULTI_AZ_TESTS}"
export ENABLE_OSS_ONLY_TESTS="${ENABLE_OSS_ONLY_TESTS}"
export ENABLE_PERSISTENT_VOLUME_TESTS="${ENABLE_PERSISTENT_VOLUME_TESTS}"
export ENABLE_POD_LOGS_TESTS="${ENABLE_POD_LOGS_TESTS}"

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

  CONFIG="${tmpfile}" ginkgo -r -progress -v "${ROOT}/git-kubo-ci/src/tests/integration-tests/"
}

main
