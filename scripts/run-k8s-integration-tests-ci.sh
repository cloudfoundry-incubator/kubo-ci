#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
cp "git-kubo-ci/specs/guestbook.yml" "${KUBO_ENVIRONMENT_DIR}/addons.yml"

git-kubo-ci/scripts/run-k8s-integration-tests.sh git-kubo-deployment "${DEPLOYMENT_NAME}" "${KUBO_ENVIRONMENT_DIR}"
