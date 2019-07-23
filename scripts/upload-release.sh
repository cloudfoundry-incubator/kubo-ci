#!/bin/bash

set -exuo pipefail

target_bosh_director() {
  local source_file
  if [[ -f source-json/source.json ]]; then
      source git-kubo-ci/scripts/set-bosh-env source-json/source.json
  else
      source git-kubo-ci/scripts/set-bosh-env source-json/metadata
  fi
}

main() {
  target_bosh_director
  bosh upload-release https://storage.googleapis.com/kubo-public/docker-35.2.3-ubuntu-xenial-315.36-20190716-163114-008878.tgz
  bosh upload-release https://storage.googleapis.com/kubo-public/docker-35.2.3-windows2019-2019.7-20190716-161813-432556.tgz
  bosh upload-release https://storage.googleapis.com/kubo-public/kubo-1.0.0-dev.102-ubuntu-xenial-456.1-20190719-232139-365998982.tgz
  bosh upload-release https://storage.googleapis.com/kubo-public/kubo-1.0.0-dev.102-windows2019-2019.4-20190719-231442-026390434.tgz
}

main
