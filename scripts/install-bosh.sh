#!/bin/sh -e

[ -z "$DEBUG" ] || set -x

. "$(dirname "$0")/lib/environment.sh"

printenv GCP_SERVICE_ACCOUNT > "$PWD/key.json"
set -x
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"

# Deploy KuBOSH
"git-kubo-deployment/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key.json"

cp "${KUBO_ENVIRONMENT_DIR}/creds.yml" "$PWD/bosh-creds/"
cp "${KUBO_ENVIRONMENT_DIR}/state.json" "$PWD/bosh-state/"

