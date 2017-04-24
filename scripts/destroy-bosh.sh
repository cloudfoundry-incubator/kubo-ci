#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

. "$(dirname "$0")/lib/environment.sh"

printenv GCP_SERVICE_ACCOUNT > "$PWD/key.json"
set -x
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"

cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}"
cp  "$PWD/gcs-bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}"

# Destroy KuBOSH
"git-kubo-deployment/bin/destroy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key.json"
