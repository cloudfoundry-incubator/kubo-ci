#!/bin/bash

set -exu -o pipefail

modify_port() {
  local property_name="${1}"
  local ops_file="$(mktemp)"

cat > "${ops_file}" <<EOF
- type: replace
  path: /${property_name}
  value: ${2}
EOF

  local temp_metadata_file="$(mktemp)"
  bosh-cli int "kubo-lock-pre/metadata" \
    -o "${ops_file}" \
    > "${temp_metadata_file}"

  cp "${temp_metadata_file}" "kubo-lock-pre/metadata"
}

main() {
  modify_port "external_kubo_port" "${PORT_NUMBER}"
  modify_port "kubernetes_master_port" "${PORT_NUMBER}"

  cp -R kubo-lock-pre/* kubo-lock
}

main
