#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

export GOPATH="${ROOT}/git-kubo-ci"
kubeconfig="gcs-kubeconfig/${KUBECONFIG_FILE}"

main() {
  if [[ ! -e "${kubeconfig}" ]]; then
    echo "Did not find kubeconfig at gcs-kubeconfig/config!"
    exit 1
  fi
  if [[ -n "$USE_SSHUTTLE" ]]; then
    bosh int kubo-lock/metadata --path=/jumpbox_ssh_key > ssh.key
    chmod 0600 ssh.key
    cidr="$(bosh int kubo-lock/metadata --path=/internal_cidr)"
    jumpbox_url="$(bosh int kubo-lock/metadata --path=/jumpbox_url)"
    sshuttle -r "jumpbox@${jumpbox_url}" "${cidr}" -e "ssh -i ssh.key -o StrictHostKeyChecking=no -o ServerAliveInterval=300 -o ServerAliveCountMax=10" --daemon
    trap 'kill -9 $(cat sshuttle.pid)' EXIT
  fi
  mkdir -p ~/.kube
  cp ${kubeconfig} ~/.kube/config

  skipped_packages=""

  if [[ "${ENABLE_MULTI_AZ_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,multiaz"
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

  ginkgo -keepGoing -r -progress -flakeAttempts=2 -skipPackage "${skipped_packages}" "${ROOT}/git-kubo-ci/src/tests/integration-tests/"
}

main
