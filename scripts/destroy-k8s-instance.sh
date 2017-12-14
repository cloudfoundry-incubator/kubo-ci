#!/bin/bash

set -exu -o pipefail

creds_path="${PWD}/gcs-bosh-creds/creds.yml"

export DEPLOYMENT_NAME="${DEPLOYMENT_NAME:-"ci-service"}"

export BOSH_CLIENT="bosh_admin"
export BOSH_CLIENT_SECRET="$(bosh-cli int "$creds_path" --path /bosh_admin_client_secret)"
export BOSH_ENVIRONMENT="$(bosh-cli int "kubo-lock/metadata" --path /internal_ip)"
export BOSH_CA_CERT="$(bosh-cli int "${creds_path}" --path=/director_ssl/ca)"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"

bosh-cli -d "${DEPLOYMENT_NAME}" -n delete-deployment
