#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_ENV="${KUBO_ENVIRONMENT_DIR}"
export BOSH_NAME=$(basename ${BOSH_ENV})
export DEBUG=1

cp "$PWD/kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"

. "git-kubo-ci/scripts/lib/deploy_utils"
. "$PWD/git-kubo-deployment/bin/set_bosh_environment"

BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET="$(get_bosh_secret)" bosh -n -e "${BOSH_ENVIRONMENT}" delete-deployment -d "tinyproxy"
