#!/bin/bash

set -eu -o pipefail

export KD=${PWD}/git-kubo-deployment
export CC_VARS_FILE=${PWD}/vm-ext-cc-vars.yml
export LOCK_FILE=${PWD}/kubo-lock/metadata
export IAAS=$(bosh int "${LOCK_FILE}" --path=/iaas)
export DEPLOYMENT_NAME="${DEPLOYMENT_NAME:-ci-service}"

create_aws_vars_file() {

cat <<EOF > "${CC_VARS_FILE}"
master_iam_instance_profile: $(bosh int "${LOCK_FILE}" --path=/master_iam_instance_profile)
worker_iam_instance_profile:  $(bosh int "${LOCK_FILE}" --path=/worker_iam_instance_profile)
cfcr_master_target_pool: $(bosh int "${LOCK_FILE}" --path=/master_target_pool)
kubernetes_cluster_tag: $(bosh int "${LOCK_FILE}" --path=/kubernetes_cluster_tag)
deployment_name: "${DEPLOYMENT_NAME}"
EOF
}

create_gcp_vars_file() {

cat <<EOF > "${CC_VARS_FILE}"
cfcr_master_service_account_address: $(bosh int "${LOCK_FILE}" --path=/service_account_master)
cfcr_worker_service_account_address: $(bosh int "${LOCK_FILE}" --path=/service_account_worker)
deployment_name: "${DEPLOYMENT_NAME}"
cfcr_backend_service: $(bosh int "${LOCK_FILE}" --path=/cfcr_backend_service)
EOF
}

create_vars_file() {
  touch "${CC_VARS_FILE}"
  if [ "${IAAS}" == "aws" ]; then
    create_aws_vars_file
  fi

  if [ "${IAAS}" == "gcp" ]; then
    create_gcp_vars_file
  fi
}


target_bosh_director() {
  BOSH_ENVIRONMENT=$(bosh int source-json/source.json --path '/target')
  BOSH_CLIENT=$(bosh int source-json/source.json --path '/client')
  BOSH_CLIENT_SECRET=$(bosh int source-json/source.json --path '/client_secret')
  BOSH_CA_CERT=$(bosh int source-json/source.json --path '/ca_cert')
  export BOSH_ENVIRONMENT BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_CA_CERT
}

update_config() {
  bosh -n update-config --name cfcr-vm-ext \
   "${KD}/manifests/cloud-config/iaas/${IAAS}/use-vm-extensions.yml" \
   --type cloud \
   --vars-file "${CC_VARS_FILE}"
}

main() {
  target_bosh_director
  create_vars_file
  update_config
}

main "$@"
