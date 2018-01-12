#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"

touch "${KUBO_ENVIRONMENT_DIR}/director-secrets.yml"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}"
cp  "$PWD/gcs-bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}"

iaas=$(bosh int kubo-lock/metadata --path=/iaas)
if [ "$iaas" = "gcp" ]; then
  set +x
  bosh int kubo-lock/metadata --path=/gcp_service_account > "$PWD/key.json"
  set -x
  "${KUBO_DEPLOYMENT_DIR}/bin/destroy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key.json"
elif [ "$iaas" = "aws" ] || [ "$iaas" = "openstack" ]; then
  set +x
  bosh int kubo-lock/metadata --path=/private_key > "$PWD/key"
  set -x
  "${KUBO_DEPLOYMENT_DIR}/bin/destroy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key"
else
  "${KUBO_DEPLOYMENT_DIR}/bin/destroy_bosh" "${KUBO_ENVIRONMENT_DIR}"
fi
