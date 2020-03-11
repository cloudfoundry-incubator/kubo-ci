#!/bin/bash

set -exu -o pipefail


export BOSH_ENVIRONMENT="$(jq -r .target source-json/source.json)"
export BOSH_CLIENT="$(jq -r .client source-json/source.json)"
export BOSH_CLIENT_SECRET="$(jq -r .client_secret source-json/source.json)"
export BOSH_CA_CERT="$(jq -r .ca_cert source-json/source.json)"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"

set +x
bosh -d "${DEPLOYMENT_NAME}" -n delete-deployment --force

export CREDHUB_SERVER="${BOSH_ENVIRONMENT/25555/8844}"
export CREDHUB_CLIENT="credhub-admin"
export CREDHUB_SECRET="$(bosh int gcs-bosh-creds/creds.yml --path=/credhub_admin_client_secret)"

tmp_uaa_ca_file="$(mktemp)"
tmp_credhub_ca_file="$(mktemp)"

trap 'rm "${tmp_uaa_ca_file}" "${tmp_credhub_ca_file}"' EXIT

bosh int "gcs-bosh-creds/creds.yml" --path="/uaa_ssl/ca" > "${tmp_uaa_ca_file}"
bosh int "gcs-bosh-creds/creds.yml" --path="/credhub_tls/ca" > "${tmp_credhub_ca_file}"

credhub login --ca-cert "${tmp_credhub_ca_file}" --ca-cert "${tmp_uaa_ca_file}"

# don't delete leading & trailing slash. This is to scope to the deployment creds we want to delete
credhub find -n "/${DEPLOYMENT_NAME}/" --output-json | jq -r .credentials[].name | grep -v managed_identity | xargs -L 1 credhub delete -n
