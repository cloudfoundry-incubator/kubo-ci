#!/bin/bash

set -exu -o pipefail

add_100_to_property() {
  local property_name="${1}"
  local property_value="$(bosh-cli int --path "/${property_name}" "kubo-lock-pre/metadata")"
  local new_property_value="$(expr "${property_value}" + 100)"
  local ops_file="$(mktemp)"

cat > "${ops_file}" <<EOF
- type: replace
  path: /${property_name}
  value: ${new_property_value}
EOF

  local temp_metadata_file="$(mktemp)"
  bosh-cli int "kubo-lock-pre/metadata" \
    -o ${ops_file} \
    > "${temp_metadata_file}"

  cp "${temp_metadata_file}" "kubo-lock-pre/metadata"
}

main() {
  add_100_to_property "external_kubo_port"
  add_100_to_property "kubernetes_master_port"

  cp -R kubo-lock-pre/* kubo-lock
}

main
