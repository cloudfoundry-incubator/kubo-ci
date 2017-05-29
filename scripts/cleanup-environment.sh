#!/bin/bash

set -eux
set -o pipefail

login_gcp() {
  lock_file="kubo-lock-repo/${POOL_NAME}/claimed/${ENV_NAME}"

  bosh-cli int "$lock_file" --path='/gcp_service_account' > gcp_service_account.json
  gcloud auth activate-service-account --key-file=gcp_service_account.json
  gcloud config set project "$(bosh-cli int - --path=/project_id < gcp_service_account.json)"
  gcloud config set compute/zone "$(bosh-cli int "$lock_file" --path='/zone')"
}

delete_tfstate() {
  if gsutil ls "gs://kubo-pipeline-store/terraform/airgap/tinyproxy/${ENV_NAME}*"; then
    gsutil rm "gs://kubo-pipeline-store/terraform/airgap/tinyproxy/${ENV_NAME}*"
  fi
}

delete_gcloud_vms() {
  lock_file="kubo-lock-repo/${POOL_NAME}/claimed/${ENV_NAME}"

  subnetwork=$(bosh-cli int "$lock_file" --path='/subnetwork')
  subnetLink=$(gcloud compute networks subnets list "$subnetwork" --format=json | bosh-cli int - --path=/0/selfLink)
  vm_names=$(gcloud  compute instances list --filter="networkInterfaces.subnetwork=$subnetLink" --format="table(name)" | tail -n +2 )

  if [ -n "${vm_names}" ]; then
    gcloud compute instances delete ${vm_names[@]} --delete-disks=all --quiet
  fi
}

login_gcp
delete_tfstate
delete_gcloud_vms
