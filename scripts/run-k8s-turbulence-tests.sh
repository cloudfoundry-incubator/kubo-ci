#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
export GOPATH="${ROOT}/git-kubo-ci"

target_bosh_director() {
  if [[ -f source-json/source.json ]]; then
    source="source-json/source.json"
  else
    source="kubo-lock/metadata"
    DEPLOYMENT_NAME="$(bosh int kubo-lock/metadata --path=/deployment_name)"
  fi
  export BOSH_DEPLOYMENT="${DEPLOYMENT_NAME}"
  source "${ROOT}/git-kubo-ci/scripts/set-bosh-env" ${source}
}

target_turbulence_api() {
  TURBULENCE_PORT=8080
  TURBULENCE_USERNAME=turbulence
  if [[ -d gcs-bosh-creds ]]; then
    TURBULENCE_HOST=$(bosh int "${ROOT}/kubo-lock/metadata" --path=/internal_ip)
    TURBULENCE_PASSWORD=$(bosh int "${ROOT}/gcs-bosh-creds/creds.yml" --path /turbulence_api_password)
    TURBULENCE_CA_CERT=$(bosh int "${ROOT}/gcs-bosh-creds/creds.yml" --path=/turbulence_api_ca/ca)
  else
    source "${ROOT}/git-kubo-ci/scripts/credhub-login" "${ROOT}/kubo-lock/metadata"
    TURBULENCE_HOST="$(bosh int "${ROOT}/kubo-lock/metadata" --path=/turbulence_api_ip)"
    cluster="$(bosh int "${ROOT}/kubo-lock/metadata" --path=/director_name)"
    TURBULENCE_PASSWORD=$(credhub get -n ${cluster}/turbulence/turbulence_api_password --quiet)
    TURBULENCE_CA_CERT=$(credhub get -n ${cluster}/turbulence/turbulence_api_ca --key ca)
  fi
  export TURBULENCE_PORT TURBULENCE_USERNAME TURBULENCE_HOST TURBULENCE_PASSWORD TURBULENCE_CA_CERT
}

create_shuttle() {
    bosh int kubo-lock/metadata --path=/jumpbox_ssh_key > ssh.key
    chmod 0600 ssh.key
    cidr="$(bosh int kubo-lock/metadata --path=/internal_cidr)"
    jumpbox_url="$(bosh int kubo-lock/metadata --path=/jumpbox_url)"
    sshuttle -r "jumpbox@${jumpbox_url}" "${cidr}" -e "ssh -i ssh.key -o StrictHostKeyChecking=no -o ServerAliveInterval=300 -o ServerAliveCountMax=10" --daemon
}

main() {
  if [[ ! -e "${KUBECONFIG_PATH}" ]]; then
    echo "Did not find kubeconfig at ${KUBECONFIG_PATH}!"
    exit 1
  fi
  mkdir -p ~/.kube/
  cp "${KUBECONFIG_PATH}" ~/.kube/config

  if bosh int kubo-lock/metadata --path /vcenter_ip &>/dev/null; then
    : #skip setting up the shuttle when testing against vpshere
  elif bosh int kubo-lock/metadata --path=/jumpbox_ssh_key &>/dev/null ; then
    create_shuttle
    trap 'kill -9 $(cat sshuttle.pid)' EXIT
  fi

  skipped_packages=""

  if [[ "${ENABLE_TURBULENCE_MASTER_FAILURE_TESTS:-false}" == "false" ]]; then
    skipped_packages="master_failure"
  fi

  if [[ "${ENABLE_TURBULENCE_WORKER_FAILURE_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,worker_failure"
  fi

  if [[ "${ENABLE_TURBULENCE_PERSISTENCE_FAILURE_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,persistence_failure"
  fi

  if [[ "${ENABLE_TURBULENCE_WORKER_DRAIN_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,worker_drain"
  fi

  target_bosh_director
  target_turbulence_api

  ginkgo -skipPackage "${skipped_packages}" -failFast -progress -r "${ROOT}/git-kubo-ci/src/tests/turbulence-tests/"
}

main
