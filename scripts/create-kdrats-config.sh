#!/usr/bin/env bash

set -euo pipefail

credhub_login() {
    local credhub_admin_secret=$(bosh int "${CREDFILE}" --path "/credhub_admin_client_secret")
    local credhub_api_url="https://$(bosh int "${VARFILE}" --path "/internal_ip"):8844"

    credhub login --client-name credhub-admin --client-secret "${credhub_admin_secret}" \
        -s "${credhub_api_url}" \
        --ca-cert <(bosh int "${CREDFILE}" --path="/credhub_tls/ca") \
        --ca-cert <(bosh int "${CREDFILE}" --path="/uaa_ssl/ca") 1>/dev/null
}

create_kdrats_config() {
    local master_host="$(bosh int ${VARFILE} --path=/kubernetes_master_host)"
    local master_port="$(bosh int ${VARFILE} --path=/kubernetes_master_port)"
    local director_name="$(bosh int ${VARFILE} --path=/director_name)"
    local ca_cert="$(bosh int <(credhub get -n "${director_name}/${BOSH_DEPLOYMENT}/tls-kubernetes" --output-json) --path=/value/ca)"
    local password="$(bosh int <(credhub get -n "${director_name}/${BOSH_DEPLOYMENT}/kubo-admin-password" --output-json) --path=/value)"

    config="$(cat "k-drats-config/$CONFIG_PATH")"
    config=$(echo "$config" | jq ".api_server_url=\"https://${master_host}:${master_port}\"")
    config=$(echo "$config" | jq ".ca_cert=\"${ca_cert}\"")
    config=$(echo "$config" | jq ".username=\"fake-kdrats-user\"")
    config=$(echo "$config" | jq ".password=\"${password}\"")

    echo "$config" > kdrats-config/config.json
}

main() {
    credhub_login
    create_kdrats_config
}

main
