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
  vms=$(gcloud  compute instances list --filter="networkInterfaces.subnetwork=$subnetLink" --format="table(name,zone)" | tail -n +2 )

  IFS=$'\n'

  for vm in $vms; do
    vm_name="$(echo $vm | awk '{print $1}')"
    vm_zone="$(echo $vm | awk '{print $2}')"

    gcloud compute instances delete "$vm_name" --zone="$vm_zone" --delete-disks=all --quiet
  done

  unset IFS
}

login_gcp
delete_gcloud_vms
