#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

metadata_path="${KUBO_ENVIRONMENT_DIR}/director.yml"
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "kubo-lock/metadata" "$metadata_path"
  touch "${KUBO_ENVIRONMENT_DIR}/director-secrets.yml"
fi

iaas=$(bosh-cli int ${metadata_path} --path=/iaas)

if [[ "$iaas" != "gcp" ]]; then
  echo "Unsupported iaas: ${iaas}"
  exit 1
fi

BOSH_EXTRA_OPS="--ops-file \"$KUBO_CI_DIR/manifests/turbulence/turbulence.yml\""
BOSH_EXTRA_OPS="${BOSH_EXTRA_OPS} --ops-file \"${KUBO_DEPLOYMENT_DIR}/bosh-deployment/jumpbox-user.yml\""

if [[ -f "$KUBO_CI_DIR/manifests/ops-files/${iaas}-cpi.yml" ]]; then
  BOSH_EXTRA_OPS="${BOSH_EXTRA_OPS} --ops-file $KUBO_CI_DIR/manifests/ops-files/${iaas}-cpi.yml"
fi

export BOSH_EXTRA_OPS

set +x
echo ${GCP_SERVICE_ACCOUNT} > "$PWD/key.json"
set -x

"${KUBO_DEPLOYMENT_DIR}/bin/deploy_bosh" "${KUBO_ENVIRONMENT_DIR}" "$PWD/key.json"

# for Concourse outputs
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  director_name=$(bosh-cli int ${metadata_path} --path=/director_name)

  echo "Storing state"
  set +x
  credhub login
  credhub set -n "/concourse/main/${director_name}/creds" -v "$(cat "${KUBO_ENVIRONMENT_DIR}/creds.yml")" -t value
  credhub set -n "/concourse/main/${director_name}/state" -v "$(cat "${KUBO_ENVIRONMENT_DIR}/state.json")" -t value
  credhub logout
  set -x
fi
