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

update_kubo() {
  # Workaround due to https://www.pivotaltracker.com/story/show/152155545
  echo "Deleting tls-kubernetes from credhub..."
  local credhub_api="https://$(bosh-cli int environment/director.yml --path=/internal_ip):8844"
  local credhub_password="$(bosh-cli int environment/creds.yml --path=/credhub_cli_password)"
  credhub login \
    -u credhub-cli \
    -p "$credhub_password" \
    -s "$credhub_api" \
    --ca-cert=<(bosh-cli int environment/creds.yml --path=/credhub_tls/ca) \
    --ca-cert=<(bosh-cli int environment/creds.yml --path=/default_ca/ca)
  credhub delete -n "$(bosh-cli int environment/director.yml --path=/director_name)/ci-service/tls-kubernetes"

  echo "Updating Kubo..."
  ${DIR}/deploy-k8s-instance.sh
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
upload_new_releases
run_upgrade_test update_kubo "$HA_MIN_SERVICE_AVAILABILITY"

# for Concourse outputs
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "${KUBO_ENVIRONMENT_DIR}/creds.yml" "$PWD/bosh-creds/"
  cp "${KUBO_ENVIRONMENT_DIR}/state.json" "$PWD/bosh-state/"
fi
