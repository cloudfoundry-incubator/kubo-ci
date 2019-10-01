#!/bin/bash

set -euo pipefail

export ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

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

  ls $RELEASE_PATH | xargs -n 1 bosh upload-release
}

main
