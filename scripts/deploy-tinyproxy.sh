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

bosh -n -e "${BOSH_ENVIRONMENT}" \
  update-cloud-config "${KUBO_DEPLOYMENT_DIR}/configurations/${IAAS}/cloud-config.yml" \
  -l "${KUBO_ENVIRONMENT_DIR}/director.yml"
bosh -n -e "${BOSH_ENVIRONMENT}" upload-stemcell "${stemcell_url}"
bosh -n -e "${BOSH_ENVIRONMENT}" deploy "${manifest_file}" -d "tinyproxy"
