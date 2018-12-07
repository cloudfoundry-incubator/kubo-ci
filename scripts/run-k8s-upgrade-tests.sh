#!/usr/bin/env bash

set -eu -o pipefail

export ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
export KUBO_ENVIRONMENT_DIR="${ROOT}/environment"
export GOPATH="${ROOT}/git-kubo-ci"

target_bosh_director() {
  BOSH_DEPLOYMENT="${DEPLOYMENT_NAME}"
  BOSH_ENVIRONMENT=$(bosh int source-json/source.json --path '/target')
  BOSH_CLIENT=$(bosh int source-json/source.json --path '/client')
  BOSH_CLIENT_SECRET=$(bosh int source-json/source.json --path '/client_secret')
  BOSH_CA_CERT=$(bosh int source-json/source.json --path '/ca_cert')
  export BOSH_DEPLOYMENT BOSH_ENVIRONMENT BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_CA_CERT
}

main() {
  target_bosh_director

  local release_tarball
  release_tarball=$(find "${ROOT}/gcs-kubo-release-tarball/" -name "*kubo-*.tgz" | head -n1)

  bosh upload-release "$release_tarball"

  if [[ "$INTERNET_ACCESS" != "false" ]]; then
    BOSH_DEPLOY_COMMAND="$ROOT/bosh-command/bosh-deploy.sh" ginkgo -r -v -progress "${ROOT}/git-kubo-ci/src/tests/upgrade-tests/"
  else
    $ROOT/bosh-command/bosh-deploy.sh
  fi
}

main
