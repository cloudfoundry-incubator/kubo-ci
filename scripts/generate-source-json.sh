#!/bin/bash

target="https://$(bosh int kubo-lock/metadata --path=/internal_ip):25555"
client="admin"
client_secret="$(bosh int gcs-bosh-creds/creds.yml --path=/admin_password)"
ca_cert="$(bosh int gcs-bosh-creds/creds.yml --path=/director_ssl/ca)"
kubernetes_master_host="$(bosh int kubo-lock/metadata --path=/kubernetes_master_host)"

jq -n \
  --arg kubernetes_master_host "$kubernetes_master_host" \
  --arg target "$target" \
  --arg client "$client" \
  --arg client_secret "$client_secret" \
  --arg ca_cert "$ca_cert" \
  '{"target": $target, "client": $client, "client_secret": $client_secret, "ca_cert": $ca_cert, "CI_kubernetes_master_host": $kubernetes_master_host}' > source-json/source.json
