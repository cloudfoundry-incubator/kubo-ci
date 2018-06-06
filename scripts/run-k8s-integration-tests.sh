#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

export GOPATH="${ROOT}/git-kubo-ci"

main() {
  source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
  setup_env "${KUBO_ENVIRONMENT_DIR}"

  local tmpfile="$(mktemp)" && echo "CONFIG=${tmpfile}"
  "${ROOT}/git-kubo-ci/scripts/generate-test-config.sh" ${KUBO_ENVIRONMENT_DIR} ${DEPLOYMENT_NAME} > "${tmpfile}"

  skipped_packages=""

  if [[ -z "${ENABLE_MULTI_AZ_TESTS+x}" ]]; then
    skipped_packages="$skipped_packages,multiaz"
  fi

  if [[ -z "${ENABLE_OSS_ONLY_TESTS+x}" ]]; then
    skipped_packages="$skipped_packages,oss_only"
  fi

  if [[ -z "${ENABLE_PERSISTENT_VOLUME_TESTS+x}" ]]; then
    skipped_packages="$skipped_packages,persistent_volume"
  fi

  if [[ -z "${ENABLE_K8S_LBS_TESTS+x}" ]]; then
    skipped_packages="$skipped_packages,k8s_lbs"
  fi

  if [[ "$skipped_packages" != "" ]]; then
    skipped_packages="$(echo $skipped_packages | cut -c 2-)"
  fi

  CONFIG="${tmpfile}" ginkgo -r -progress -v -skipPackage "${skipped_packages}" "${ROOT}/git-kubo-ci/src/tests/integration-tests/"
}

main
