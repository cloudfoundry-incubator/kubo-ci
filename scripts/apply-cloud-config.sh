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

update_runtime_config() {
  bosh -n update-runtime-config --name=tinyproxy <(echo "$RUNTIME_CONFIG_YML")
}

source_json="$1"

main() {
  local tmp_cloud_config
  [ ! -f "$source_json" ] && usage

  target_bosh_director

  tmp_cloud_config=$(mktemp)

  bosh cloud-config > "$tmp_cloud_config"
  bosh -n update-cloud-config "$tmp_cloud_config" \
   -o "${PWD}/git-kubo-ci/manifests/ops-files/vsphere-proxy-cloud-config.yml" \
   -v "proxy_static_ip=${PROXY_STATIC_IP}"
}

main "$@"
