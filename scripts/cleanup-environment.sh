#!/bin/bash

set -exu -o pipefail

login_gcp() {
  if bosh int "${ENV_FILE}" --path='/gcp_service_account' &> /dev/null; then
    bosh int "${ENV_FILE}" --path='/gcp_service_account' > gcp_service_account.json
  elif [[ -n "${GCP_SERVICE_ACCOUNT}" ]]; then
    set +x
    echo "${GCP_SERVICE_ACCOUNT}" > gcp_service_account.json
    set -x
  fi
  gcloud auth activate-service-account --key-file=gcp_service_account.json
  gcloud config set project "$(bosh int - --path=/project_id < gcp_service_account.json)"
  gcloud config set compute/zone "$(bosh int "${ENV_FILE}" --path='/zone')"
}

delete_gcloud_vms() {
  subnetwork=$(bosh int "${ENV_FILE}" --path='/subnetwork')
  subnetLink=$(gcloud compute networks subnets list "$subnetwork" --format=json | bosh int - --path=/0/selfLink)
  vm_names=$(gcloud  compute instances list --filter="networkInterfaces.subnetwork=$subnetLink" --format="table(name)" | tail -n +2 )

  if [ -n "${vm_names}" ]; then
    gcloud compute instances delete ${vm_names[@]} --delete-disks=all --quiet
  fi
}

login_gcp
delete_gcloud_vms
