#!/bin/bash

set -eu -o pipefail

export KD=${PWD}/git-kubo-deployment
export CC_VARS_FILE=${PWD}/vm-ext-cc-vars.yml
export LOCK_FILE=${PWD}/kubo-lock/metadata

create_vars_file() {

cat <<EOF > ${CC_VARS_FILE}
master_iam_instance_profile: $(bosh int ${LOCK_FILE} --path=/master_iam_instance_profile)
worker_iam_instance_profile:  $(bosh int ${LOCK_FILE} --path=/worker_iam_instance_profile)
cfcr_master_target_pool: $(bosh int ${LOCK_FILE} --path=/master_target_pool)
kubernetes_cluster_tag: $(bosh int ${LOCK_FILE} --path=/kubernetes_cluster_tag)
deployment_name: ci-service
EOF

}

update_config() {
  local iaas=$(bosh int ${LOCK_FILE} --path=/iaas)

  bosh -n update-config --name cfcr-vm-ext \
   ${KD}/manifests/cloud-config/iaas/${iaas}/use-vm-extensions.yml \
   --type cloud \
   --vars-file ${CC_VARS_FILE}
}

main() {
  create_vars_file
  update_config
}

main "$@"
