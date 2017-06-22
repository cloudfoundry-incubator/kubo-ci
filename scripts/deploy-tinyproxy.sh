#!/bin/bash

set -exu -o pipefail

source "$(dirname "$0")/lib/environment.sh"

export BOSH_ENV="${KUBO_ENVIRONMENT_DIR}"
export BOSH_NAME=$(basename ${BOSH_ENV})
export DEBUG=1

cp "$PWD/kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"

source "$PWD/git-kubo-deployment/bin/lib/deploy_utils"
source "$PWD/git-kubo-deployment/bin/set_bosh_environment"

cloud_config_file="$PWD/git-kubo-ci/utils/tinyproxy/cloud-config.yml"
manifest_file="$PWD/git-kubo-ci/utils/tinyproxy/manifest.yml"
stemcell_url=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path='/stemcell_url')

BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET="$(get_bosh_secret)" bosh-cli -n -e "${BOSH_ENVIRONMENT}" update-cloud-config "${KUBO_DEPLOYMENT_DIR}/configurations/${IAAS}/cloud-config.yml" -l "${KUBO_ENVIRONMENT_DIR}/director.yml"
BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET="$(get_bosh_secret)" bosh-cli -n -e "${BOSH_ENVIRONMENT}" upload-stemcell "$stemcell_url"
BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET="$(get_bosh_secret)" bosh-cli -n -e "${BOSH_ENVIRONMENT}" deploy "$manifest_file" -d "tinyproxy"
