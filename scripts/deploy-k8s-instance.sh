#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

tarball_name=$(ls $PWD/gcs-kubo-release-tarball/kubo-*.tgz | head -n1)

cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

cp "$tarball_name" "git-kubo-deployment/../kubo-release.tgz"

"git-kubo-deployment/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"

iaas=$(bosh-cli int ${KUBO_ENVIRONMENT_DIR}/director.yml --path=/iaas)
iaas_cc_opsfile="${PWD}/git-kubo-ci/manifests/ops-files/${iaas}-k8s-cloud-config.yml"
static_network_cc_opsfile="${PWD}/git-kubo-ci/manifests/ops-files/${iaas}-static-network-cloud-config.yml"

if [[ -n "${CC_STATIC_NETWORK}" && "${CC_STATIC_NETWORK}" = true && -f "$static_network_cc_opsfile" ]]; then
  export CLOUD_CONFIG_OPS_FILE="${static_network_cc_opsfile}"
elif [[ -f "$iaas_cc_opsfile" ]]; then
  export CLOUD_CONFIG_OPS_FILE="${iaas_cc_opsfile}"
fi

"git-kubo-deployment/bin/deploy_k8s" "${KUBO_ENVIRONMENT_DIR}" ci-service local
