#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"

HA_MIN_SERVICE_AVAILABILITY="${HA_MIN_SERVICE_AVAILABILITY:-1}"

update_bosh() {
  echo "Updating BOSH..."
  ${DIR}/install-bosh.sh
}

# copy state and creds so that deploy_bosh has the correct context
copy_state_and_creds() {
  cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "$PWD/gcs-bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}/"
  cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
  touch "${KUBO_ENVIRONMENT_DIR}/director-secrets.yml"
}

if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  copy_state_and_creds
fi
run_upgrade_test update_bosh $HA_MIN_SERVICE_AVAILABILITY

# for Concourse outputs
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "${KUBO_ENVIRONMENT_DIR}/creds.yml" "$PWD/bosh-creds/"
  cp "${KUBO_ENVIRONMENT_DIR}/state.json" "$PWD/bosh-state/"
fi
