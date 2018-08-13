#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

export GOPATH="${ROOT}/git-kubo-ci"

main() {
  if [[ ! -e gcs-kubeconfig/config ]]; then
    echo "Did not find kubeconfig at gcs-kubeconfig/config!"
    exit 1
  fi
  mkdir -p ~/.kube
  cp gcs-kubeconfig/config ~/.kube/config

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

  ginkgo -keepGoing -r -progress -skipPackage "${skipped_packages}" "${ROOT}/git-kubo-ci/src/tests/integration-tests/"
}

main
