#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"
. "$DIR/lib/utils.sh"

update_bosh() {
  echo "Updating BOSH..."
  ${DIR}/install-bosh.sh
}

update_kubo() {
  echo "Updating Kubo..."
  export DEPLOYMENT_NAME=ci-service
  ${DIR}/deploy-k8s-instance.sh
}

# copy state and creds so that deploy_bosh has the correct context
copy_state_and_creds() {
  cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "$PWD/gcs-bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}/"
  cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
  touch "${KUBO_ENVIRONMENT_DIR}/director-secrets.yml"
}

tmpfile=$(mktemp)
$DIR/generate-test-config.sh "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}" > "${tmpfile}"
export CONFIG="${tmpfile}"

export GOPATH="${DIR}/.."
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  copy_state_and_creds
fi

set_kubeconfig
ginkgo -progress -v "$DIR/../src/tests/upgrade-tests"

# for Concourse outputs
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "${KUBO_ENVIRONMENT_DIR}/creds.yml" "$PWD/bosh-creds/"
  cp "${KUBO_ENVIRONMENT_DIR}/state.json" "$PWD/bosh-state/"
fi
