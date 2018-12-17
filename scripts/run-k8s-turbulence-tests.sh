#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

kubeconfig="gcs-kubeconfig/${KUBECONFIG_FILE}"
export GOPATH="${ROOT}/git-kubo-ci"

target_bosh_director() {
  source="source-json/source.json"
  if [[ ! -f source-json/source.json ]]; then
    source="kubo-lock/metadata"
    DEPLOYMENT_NAME="$(bosh int kubo-locks/metadata --path=/deployment_name)"
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
  TURBULENCE_PASSWORD=$(bosh int "${ROOT}/gcs-bosh-creds/creds.yml" --path /turbulence_api_password)
  TURBULENCE_CA_CERT=$(bosh int "${ROOT}/gcs-bosh-creds/creds.yml" --path=/turbulence_api_ca/ca)
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

  ginkgo -skipPackage "${skipped_packages}" -failFast -progress -r "${ROOT}/git-kubo-ci/src/tests/turbulence-tests/"
}

main
