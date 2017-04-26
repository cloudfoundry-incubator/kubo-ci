#!/bin/sh -e

. "$(dirname "$0")/lib/environment.sh"
bosh-cli int kubo-lock/metadata --path=/gcp_service_account > "$PWD/key.json"

set -x
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"


cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}"
cp  "$PWD/gcs-bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}"

"${KUBO_DEPLOYMENT_LOCATION}/bin/destroy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key.json"
