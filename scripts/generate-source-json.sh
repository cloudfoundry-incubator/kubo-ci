#!/bin/bash

set -euo pipefail

target="https://$(bosh int kubo-lock/metadata --path=/internal_ip):25555"
client="admin"
client_secret="$(bosh int gcs-bosh-creds/creds.yml --path=/admin_password)"
ca_cert="$(bosh int gcs-bosh-creds/creds.yml --path=/director_ssl/ca)"

jq -n \
  --arg target "$target" \
  --arg client "$client" \
  --arg client_secret "$client_secret" \
  --arg ca_cert "$ca_cert" \
  '{"target": $target, "client": $client, "client_secret": $client_secret, "ca_cert": $ca_cert}' > source-json/source.json
