#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

metadata_path="${KUBO_ENVIRONMENT_DIR}/director.yml"
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock-pre/metadata" "$metadata_path"

set +x
director_file_path="${KUBO_ENVIRONMENT_DIR}/director.yml"
creds_file_path="${KUBO_ENVIRONMENT_DIR}/creds.yml"
credhub_admin_secret=$(bosh int $creds_file_path --path=/credhub_admin_client_secret)
credhub_api_url="https://$(bosh int $director_file_path --path=/internal_ip):8844"
director_name=$(bosh int $director_file_path --path=/director_name)
bosh_admin_client_secret=$(bosh int "${creds_file_path}" --path=/bosh_admin_client_secret)

tmp_uaa_ca_file="$(mktemp)"
tmp_credhub_ca_file="$(mktemp)"
trap 'rm "${tmp_uaa_ca_file}" "${tmp_credhub_ca_file}"' EXIT

bosh int "$creds_file_path" --path="/uaa_ssl/ca" > "${tmp_uaa_ca_file}"
bosh int "$creds_file_path" --path="/credhub_tls/ca" > "${tmp_credhub_ca_file}"

credhub login --client-name credhub-admin --client-secret "${credhub_admin_secret}" -s "${credhub_api_url}" --ca-cert "${tmp_credhub_ca_file}" --ca-cert "${tmp_uaa_ca_file}"

set -x
DEPLOYMENT_NAME=${DEPLOYMENT_NAME:-"ci-service"}
credhub delete -n "/${director_name}/${DEPLOYMENT_NAME}/tls-kubernetes"

"$KUBO_DEPLOYMENT_DIR/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"
master_ip=$(bosh int <(BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET="${bosh_admin_client_secret}" bosh -e environment -d "${DEPLOYMENT_NAME}" vms --json) --path=/Tables/0/Rows/0/ips)

cp -R kubo-lock-pre/* kubo-lock
echo "master_ip: ${master_ip}" >> kubo-lock/metadata
