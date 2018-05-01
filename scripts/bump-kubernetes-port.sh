#!/bin/bash

set -exu -o pipefail

bump_port() {
  local property_name="${1}"
  local bump_amount="${2}"
  local ops_file="$(mktemp)"
  local current_port=$(bosh int "kubo-lock-pre/metadata" --path=/${property_name})


cat > "${ops_file}" <<EOF
- type: replace
  path: /${property_name}
  value: $(( current_port + bump_amount ))
EOF

  local temp_metadata_file="$(mktemp)"
  bosh-cli int "kubo-lock-pre/metadata" \
    -o "${ops_file}" \
    > "${temp_metadata_file}"

  cp "${temp_metadata_file}" "kubo-lock-pre/metadata"
}

main() {
  # Uncomment this for CF routing
  #bump_port "external_kubo_port" "${BUMP_AMOUNT}"
  bump_port "kubernetes_master_port" "${BUMP_AMOUNT}"

  cp -R kubo-lock-pre/* kubo-lock
}

main
