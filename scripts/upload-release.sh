#!/bin/bash

set -euo pipefail

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

  bosh upload-release "$RELEASES_PATH"
}

