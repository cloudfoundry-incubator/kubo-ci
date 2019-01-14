#!/usr/bin/env bash

set -eu -o pipefail

export ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
export GOPATH="${ROOT}/git-kubo-ci"

target_bosh_director() {
  if [[ -f source-json/source.json ]]; then
    source="source-json/source.json"
  else
    source="kubo-lock/metadata"
    DEPLOYMENT_NAME="$(bosh int kubo-lock/metadata --path=/deployment_name)"
  fi
  export BOSH_DEPLOYMENT="${DEPLOYMENT_NAME}"
  source "${ROOT}/git-kubo-ci/scripts/set-bosh-env" ${source}
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
