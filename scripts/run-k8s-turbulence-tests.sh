#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

kubeconfig="gcs-kubeconfig/${KUBECONFIG_FILE}"
export GOPATH="${ROOT}/git-kubo-ci"

target_bosh_director() {
  if [[ -f source-json/source.json ]]; then
    source="source-json/source.json"
  else
    source="kubo-lock/metadata"
    DEPLOYMENT_NAME="$(bosh int kubo-lock/metadata --path=/deployment_name)"
  fi
  BOSH_DEPLOYMENT="${DEPLOYMENT_NAME}"
  BOSH_ENVIRONMENT=$(bosh int "${source}" --path '/target')
  BOSH_CLIENT=$(bosh int "${source}" --path '/client')
  BOSH_CLIENT_SECRET=$(bosh int "${source}" --path '/client_secret')
  BOSH_CA_CERT=$(bosh int "${source}" --path '/ca_cert')
  export BOSH_DEPLOYMENT BOSH_ENVIRONMENT BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_CA_CERT
}

target_turbulence_api() {
  TURBULENCE_PORT=8080
  TURBULENCE_USERNAME=turbulence
  TURBULENCE_HOST=$(bosh int "${ROOT}/kubo-lock/metadata" --path=/internal_ip)
  if [[ -f turbulence/creds.yml  ]]; then
    TURBULENCE_PASSWORD=$(bosh int "${ROOT}/turbulence/creds.yml" --path /turbulence_api_password)
    TURBULENCE_CA_CERT=$(bosh int "${ROOT}/turbulence/creds.yml" --path=/turbulence_api_ca/ca)
  else
    TURBULENCE_PASSWORD=$(bosh int "${ROOT}/gcs-bosh-creds/creds.yml" --path /turbulence_api_password)
    TURBULENCE_CA_CERT=$(bosh int "${ROOT}/gcs-bosh-creds/creds.yml" --path=/turbulence_api_ca/ca)
  fi
  export TURBULENCE_PORT TURBULENCE_USERNAME TURBULENCE_HOST TURBULENCE_PASSWORD TURBULENCE_CA_CERT
}

main() {
  if [[ ! -e "${kubeconfig}" ]]; then
    echo "Did not find kubeconfig at gcs-kubeconfig/${KUBECONFIG_FILE}!"
    exit 1
  fi

  mkdir -p ~/.kube
  cp ${kubeconfig} ~/.kube/config

  skipped_packages=""

  if [[ "${ENABLE_TURBULENCE_MASTER_FAILURE_TESTS:-false}" == "false" ]]; then
    skipped_packages="$skipped_packages,master_failure"
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

  if bosh int kubo-lock/metadata --path=/jumpbox_ssh_key &>/dev/null ; then
    bosh int kubo-lock/metadata --path=/jumpbox_ssh_key > ssh.key
    chmod 0600 ssh.key
    cidr="$(bosh int kubo-lock/metadata --path=/internal_cidr)"
    jumpbox_url="$(bosh int kubo-lock/metadata --path=/jumpbox_url)"
    sshuttle -r "jumpbox@${jumpbox_url}" "${cidr}" -e "ssh -i ssh.key -o StrictHostKeyChecking=no -o ServerAliveInterval=300 -o ServerAliveCountMax=10" --daemon
    trap 'kill -9 $(cat sshuttle.pid)' EXIT
  fi

  ginkgo -skipPackage "${skipped_packages}" -failFast -progress -r "${ROOT}/git-kubo-ci/src/tests/turbulence-tests/"
}

main
