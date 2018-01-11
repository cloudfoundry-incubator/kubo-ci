#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

metadata_path="${KUBO_ENVIRONMENT_DIR}/director.yml"
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "kubo-lock/metadata" "$metadata_path"

  if [ -f "gcs-bosh-creds/creds.yml" ]; then
    cp gcs-bosh-creds/creds.yml "${KUBO_ENVIRONMENT_DIR}/creds.yml"
  fi
  if [ -f "gcs-bosh-state/state.json" ]; then
    cp gcs-bosh-state/state.json "${KUBO_ENVIRONMENT_DIR}/state.json"
  fi

  touch "${KUBO_ENVIRONMENT_DIR}/director-secrets.yml"
fi

iaas=$(bosh-cli int $metadata_path --path=/iaas)

BOSH_EXTRA_OPS=""
# This means USE_TURBULENCE is set and not blank #Bashisms
if [[ ! -z ${USE_TURBULENCE+x} ]] && [[ ! -z "${USE_TURBULENCE}" ]]; then
  BOSH_EXTRA_OPS="--ops-file \"${KUBO_DEPLOYMENT_DIR}/bosh-deployment/turbulence.yml\""
fi

BOSH_EXTRA_OPS="${BOSH_EXTRA_OPS} --ops-file \"${KUBO_DEPLOYMENT_DIR}/bosh-deployment/jumpbox-user.yml\""

if [[ -f "$KUBO_CI_DIR/manifests/ops-files/${iaas}-cpi.yml" ]]; then
  BOSH_EXTRA_OPS="${BOSH_EXTRA_OPS} --ops-file $KUBO_CI_DIR/manifests/ops-files/${iaas}-cpi.yml"
fi

export BOSH_EXTRA_OPS

set +x
echo "Deploying BOSH"

if [ "$iaas" = "gcp" ]; then
  if [[ ! -z "${GCP_SERVICE_ACCOUNT+x}" ]] && [[ "$GCP_SERVICE_ACCOUNT" != "" ]]; then
    echo "$GCP_SERVICE_ACCOUNT" >> "$PWD/key.json"
  else
    bosh-cli int $metadata_path --path=/gcp_service_account > "$PWD/key.json"
  fi
  "${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key.json"
elif [ "$iaas" = "aws" ]; then
  bosh-cli int $metadata_path --path=/private_key > "$PWD/key"
  "${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key"
elif [ "$iaas" = "openstack" ]; then
  bosh-cli int $metadata_path --path=/private_key > "$PWD/key"
  "${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key"
else
  "${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}"
fi

# for Concourse outputs
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "${KUBO_ENVIRONMENT_DIR}/creds.yml" "$PWD/bosh-creds/"
  cp "${KUBO_ENVIRONMENT_DIR}/state.json" "$PWD/bosh-state/"
fi
