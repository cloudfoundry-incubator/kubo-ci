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

