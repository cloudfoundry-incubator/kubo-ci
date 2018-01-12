#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1
BOSH_ENV="${KUBO_ENVIRONMENT_DIR}"

cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

source "git-kubo-deployment/bin/set_bosh_environment"

turbulence_release_url=$(bosh int "git-kubo-ci/manifests/turbulence/runtime-config.yml" --path='/releases/name=turbulence/url')
bosh -n -e "${BOSH_ENVIRONMENT}" upload-release "$turbulence_release_url"

runtime_config=$(bosh int "git-kubo-ci/manifests/turbulence/runtime-config.yml" --vars-file ${KUBO_ENVIRONMENT_DIR}/director.yml --vars-file ${KUBO_ENVIRONMENT_DIR}/creds.yml)

echo "$runtime_config"  |  bosh -n -e "${BOSH_ENVIRONMENT}" update-runtime-config --name=turbulence -
