#!/bin/bash

set -exu -o pipefail

usage() {
  echo "USAGE: $0 [source-json]"
  echo
  exit 1
}

target_bosh_director() {
  export BOSH_ENVIRONMENT=$(bosh int $source_json --path '/target')
  export BOSH_CLIENT=$(bosh int $source_json --path '/client')
  export BOSH_CLIENT_SECRET=$(bosh int $source_json --path '/client_secret')
  export BOSH_CA_CERT=$(bosh int $source_json --path '/ca_cert')
}

source_json="$1"

main() {
  local tmp_cloud_config
  [ ! -f "$source_json" ] && usage

  target_bosh_director

  tmp_cloud_config=$(mktemp)

  bosh cloud-config > "$tmp_cloud_config"
  bosh -n update-cloud-config "$tmp_cloud_config" "${OPS}"
}

main "$@"
