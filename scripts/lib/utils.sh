#!/usr/bin/env bash

set -eu
set -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)

set_variables() {
  bosh_name="$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/director_name)"
  cluster_name="${bosh_name}/${DEPLOYMENT_NAME}"
  host="$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/kubernetes_master_host)"
  port="$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/kubernetes_master_port)"
  api_url="https://${host}:${port}"

  echo "cluster_name=$cluster_name" "api_url=$api_url"
}

setup_env() {
  KUBO_ENVIRONMENT_DIR="${1}"
  mkdir -p "${KUBO_ENVIRONMENT_DIR}"
  cp "${ROOT}/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "${ROOT}/kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

  "${ROOT}/git-kubo-deployment/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"
  if [[ -f "${ROOT}/git-kubo-deployment/bin/credhub_login" ]]; then
    "${ROOT}/git-kubo-deployment/bin/credhub_login" "${KUBO_ENVIRONMENT_DIR}"
    eval "$(set_variables)"
    "${ROOT}/git-kubo-deployment/bin/set_kubeconfig" "${cluster_name}" "${api_url}"
  else
    "${ROOT}/git-kubo-deployment/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}"
  fi
}

