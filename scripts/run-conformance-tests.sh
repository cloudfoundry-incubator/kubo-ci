#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/run_test_suite.sh"

copy_state_and_creds() {
  cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
  "$PWD/git-kubo-deployment/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "ci-service"
}

if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  copy_state_and_creds
fi

if [ -z ${CONFORMANCE_RESULTS_DIR+x} ]; then
  echo "Error: CONFORMANCE_RESULTS_DIR is not set, exiting..."
  exit 1
fi

GOPATH="$KUBO_CI_DIR"
export GOPATH
export PATH_TO_KUBECONFIG="$HOME/.kube/config"
export CONFORMANCE_RESULTS_DIR="$PWD/$CONFORMANCE_RESULTS_DIR"
export RELEASE_TARBALL="$PWD/$KUBO_DEPLOYMENT_DIR/kubo-release.tgz"

kubo::tests::run_test_suite "$GOPATH/src/tests/conformance"
