#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

export GOPATH="${ROOT}/git-kubo-ci"

main() {
  local tmpfile

  # shellcheck source=lib/utils.sh
  source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
  setup_env "${KUBO_ENVIRONMENT_DIR}"

  tmpfile="$(mktemp)" && echo "CONFIG=${tmpfile}"
  "${ROOT}/git-kubo-ci/scripts/generate-test-config.sh" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}" > "${tmpfile}"

  skipped_packages=""

  if [[ "${ENABLE_MULTI_AZ_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,multiaz"
  fi

  if [[ "${ENABLE_OSS_ONLY_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,oss_only"
  fi

  if [[ "${ENABLE_PERSISTENT_VOLUME_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,persistent_volume"
  fi

  if [[ "${ENABLE_K8S_LBS_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,k8s_lbs"
  fi

  if [[ "${ENABLE_CIDR_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,cidrs"
  fi

  if [[ "$skipped_packages" != "" ]]; then
    skipped_packages="$(echo $skipped_packages | cut -c 2-)"
  fi

  CONFIG="${tmpfile}" ginkgo -keepGoing -r -progress -skipPackage "${skipped_packages}" "${ROOT}/git-kubo-ci/src/tests/integration-tests/"
}

main
