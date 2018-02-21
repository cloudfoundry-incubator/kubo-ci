#!/bin/bash

set -exu -o pipefail

source "$(dirname "$0")/lib/environment.sh"

BOSH_ENV="${KUBO_ENVIRONMENT_DIR}"
BOSH_NAME="$(basename "${BOSH_ENV}")"
DEBUG=1
export BOSH_ENV BOSH_NAME DEBUG

cp "$PWD/kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"

source "$PWD/git-kubo-deployment/bin/lib/deploy_utils"
source "$PWD/git-kubo-deployment/bin/set_bosh_environment"

manifest_file="$PWD/git-kubo-ci/manifests/tinyproxy/manifest.yml"
stemcell_url=$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path='/stemcell_url')

BOSH_CLIENT=bosh_admin
BOSH_CLIENT_SECRET="$(get_bosh_secret)"
export BOSH_CLIENT BOSH_CLIENT_SECRET

bosh -n update-cloud-config "${KUBO_DEPLOYMENT_DIR}/configurations/${IAAS}/cloud-config.yml" \
  -l "${KUBO_ENVIRONMENT_DIR}/director.yml" -o "${PWD}/git-kubo-ci/manifests/ops-files/vsphere-proxy-cloud-config.yml" -v "proxy_static_ip=${PROXY_STATIC_IP}"

bosh -n upload-stemcell "${stemcell_url}"

opsfile_arguments=""
if [[ -n ${PROXY_STATIC_IP} ]]; then
  opsfile_arguments=" -o ${PWD}/git-kubo-ci/manifests/ops-files/airgap-tinyproxy.yml -v proxy_static_ip=${PROXY_STATIC_IP}"
fi

bosh -n -d tinyproxy deploy "${manifest_file}" ${opsfile_arguments}
