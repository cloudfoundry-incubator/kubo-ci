#!/bin/bash

set -exu -o pipefail

creds_path="${PWD}/gcs-bosh-creds/creds.yml"

export BOSH_CLIENT="bosh_admin"
export BOSH_CLIENT_SECRET="$(bosh int "$creds_path" --path /bosh_admin_client_secret)"
export BOSH_ENVIRONMENT="$(bosh int "kubo-lock/metadata" --path /internal_ip)"
export BOSH_CA_CERT="$(bosh int "${creds_path}" --path=/director_ssl/ca)"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"

set +x
bosh -d "${DEPLOYMENT_NAME}" -n delete-deployment
