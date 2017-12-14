#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

calculate_port_num() {
  offset=$(printf "%d" "'${director_name}")
  port=$(( ${CFCR_ROUTING_PORT_RANGE_START} + ${offset} ))
  echo -n "${port}"
}

environment_dir="environment"
mkdir -p ${environment_dir}

metadata_path="kubo-lock/metadata"
director_name=$(bosh-cli int ${metadata_path} --path=/director_name)

echo "Building envrionment"

director_config="${environment_dir}/director.yml"

cp "${metadata_path}" "${director_config}"
echo "${CFCR_GENERAL}" >> "${director_config}"
echo "${CFCR_IAAS}" >> "${director_config}"
echo "${CFCR_ROUTING}" >> "${director_config}"

if [ -n "${CFCR_ROUTING_PORT_RANGE_START}" ]; then
  echo "kubernetes_master_port: $(calculate_port_num)" >> "${director_config}"
fi

echo "Getting creds"

credhub login
set +x

credhub get -n "/concourse/main/${director_name}/creds" --output-json | jq -r .value > bosh-creds/creds.yml

set -x

credhub logout
