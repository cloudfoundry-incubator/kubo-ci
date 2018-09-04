#!/bin/bash

set -eu -o pipefail

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="${ROOT}/bosh.log"
export DEBUG=0

metadata_path="${KUBO_ENVIRONMENT_DIR}/director.yml"
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "${ROOT}/kubo-lock/metadata" "$metadata_path"

  if [ -f "${ROOT}/gcs-bosh-creds/creds.yml" ]; then
    cp "${ROOT}/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/creds.yml"
  fi
  if [ -f "${ROOT}/gcs-bosh-state/state.json" ]; then
    cp "${ROOT}/gcs-bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}/state.json"
  fi

  touch "${KUBO_ENVIRONMENT_DIR}/director-secrets.yml"
fi

iaas=$(bosh int $metadata_path --path=/iaas)

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

iaas=$(bosh int $metadata_path --path=/iaas)
iaas_cc_opsfile="$KUBO_CI_DIR/manifests/ops-files/${iaas}-k8s-cloud-config.yml"

CLOUD_CONFIG_OPS_FILE=${CLOUD_CONFIG_OPS_FILE:-""}
if [[ -f "$KUBO_CI_DIR/manifests/ops-files/$CLOUD_CONFIG_OPS_FILE" ]]; then
  CLOUD_CONFIG_OPS_FILES="$KUBO_CI_DIR/manifests/ops-files/$CLOUD_CONFIG_OPS_FILE"
elif [[ -f "$iaas_cc_opsfile" ]]; then
  CLOUD_CONFIG_OPS_FILES="${iaas_cc_opsfile}"
fi
export CLOUD_CONFIG_OPS_FILES

echo "Deploying BOSH"

if [ "$iaas" = "gcp" ]; then
  if [[ ! -z "${GCP_SERVICE_ACCOUNT+x}" ]] && [[ "$GCP_SERVICE_ACCOUNT" != "" ]]; then
    echo "$GCP_SERVICE_ACCOUNT" >> "${ROOT}/key.json"
  else
    bosh int $metadata_path --path=/gcp_service_account > "${ROOT}/key.json"
  fi
  "${ROOT}/git-kubo-ci/scripts/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "${ROOT}/key.json"
elif [ "$iaas" = "aws" ]; then
  bosh int $metadata_path --path=/private_key > "${ROOT}/key"
  "${ROOT}/git-kubo-ci/scripts/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "${ROOT}/key"
elif [ "$iaas" = "openstack" ]; then
  bosh int $metadata_path --path=/private_key > "${ROOT}/key"
  "${ROOT}/git-kubo-ci/scripts/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "${ROOT}/key"
else
  "${ROOT}/git-kubo-ci/scripts/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}"
fi

"${ROOT}/git-kubo-ci/scripts/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"

# for Concourse outputs
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "${KUBO_ENVIRONMENT_DIR}/creds.yml" "${ROOT}/bosh-creds/"
  cp "${KUBO_ENVIRONMENT_DIR}/state.json" "${ROOT}/bosh-state/"
fi
