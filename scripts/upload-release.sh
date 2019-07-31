#!/bin/bash

set -exuo pipefail

: "${RELEASE_LIST:?Need release list}"

target_bosh_director() {
  local source_file
  if [[ -f source-json/source.json ]]; then
      source git-kubo-ci/scripts/set-bosh-env source-json/source.json
  else
      source git-kubo-ci/scripts/set-bosh-env source-json/metadata
  fi
}

upload_releases(){
  releases="$RELEASE_LIST"
  for release in $RELEASE_LIST
  do
    bosh upload-release $(find "${release}" -name *.tgz)
  done
}

main() {
  target_bosh_director

  upload_releases
}

main
