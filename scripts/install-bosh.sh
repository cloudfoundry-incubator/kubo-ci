#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

# This means DO_UPGRADE is set and not blank #Bashisms
if [[ ! -z ${DO_UPGRADE+x} ]] && [[ ! -z "${DO_UPGRADE}" ]]; then
  cp "$PWD/bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}"
  cp "$PWD/bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}" 
fi

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
touch "${KUBO_ENVIRONMENT_DIR}/director-secrets.yml"

iaas=$(bosh-cli int kubo-lock/metadata --path=/iaas)

BOSH_EXTRA_OPS=""
# This means USE_TURBULENCE is set and not blank #Bashisms
if [[ ! -z ${USE_TURBULENCE+x} ]] && [[ ! -z "${USE_TURBULENCE}" ]]; then
  BOSH_EXTRA_OPS="--ops-file \"git-kubo-ci/manifests/turbulence/turbulence.yml\""
fi

BOSH_EXTRA_OPS="${BOSH_EXTRA_OPS} --ops-file \"${KUBO_DEPLOYMENT_DIR}/bosh-deployment/jumpbox-user.yml\""

if [[ -f "${PWD}/git-kubo-ci/manifests/ops-files/${iaas}-cpi.yml" ]]; then
  BOSH_EXTRA_OPS="${BOSH_EXTRA_OPS} --ops-file ${PWD}/git-kubo-ci/manifests/ops-files/${iaas}-cpi.yml"
fi

export BOSH_EXTRA_OPS

if [ "$iaas" = "gcp" ]; then
  set +x
  bosh-cli int kubo-lock/metadata --path=/gcp_service_account > "$PWD/key.json"
  set -x
  "${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key.json"
elif [ "$iaas" = "aws" ]; then
  set +x
  bosh-cli int kubo-lock/metadata --path=/private_key > "$PWD/key"
  set -x
  "${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key"
elif [ "$iaas" = "openstack" ]; then
  set +x
  bosh-cli int kubo-lock/metadata --path=/private_key > "$PWD/key"
  set -x
  "${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key"
else
  "${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}"
fi

cp "${KUBO_ENVIRONMENT_DIR}/creds.yml" "$PWD/bosh-creds/"
cp "${KUBO_ENVIRONMENT_DIR}/state.json" "$PWD/bosh-state/"
