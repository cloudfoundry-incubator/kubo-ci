#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"

credhub_login() {
  credhub_admin_secret=$(bosh int "${CREDFILE}" --path "/credhub_admin_client_secret")
  credhub_api_url="https://$(bosh int "${VARFILE}" --path "/internal_ip"):8844"
  credhub login --client-name credhub-admin --client-secret "${credhub_admin_secret}" \
    -s "${credhub_api_url}" \
    --ca-cert <(bosh int "${CREDFILE}" --path="/credhub_tls/ca") \
    --ca-cert <(bosh int "${CREDFILE}" --path="/uaa_ssl/ca") 1>/dev/null
}

generate_test_config() {
  local environment="$1"
  local deployment="$2"
  local enable_multi_az_tests="${ENABLE_MULTI_AZ_TESTS:-false}"
  local enable_turbulence_worker_drain_tests="${ENABLE_TURBULENCE_WORKER_DRAIN_TESTS:-false}"
  local enable_turbulence_worker_failure_tests="${ENABLE_TURBULENCE_WORKER_FAILURE_TESTS:-false}"
  local enable_turbulence_master_failure_tests="${ENABLE_TURBULENCE_MASTER_FAILURE_TESTS:-false}"
  local enable_turbulence_persistence_failure_tests="${ENABLE_TURBULENCE_PERSISTENCE_FAILURE_TESTS:-false}"
  local cidr_vars_file="${CIDR_VARS_FILE:-git-kubo-ci/manifests/vars-files/default-cidrs.yml}"

  shift 2
  for arg in "$@"; do
    local flag="${arg%%=*}"
    local value="${arg##*=}"
    case "$flag" in
      --enable-multi-az-tests)
	enable_multi_az_tests=true
	;;
      --enable-turbulence-worker-drain-tests)
        enable_turbulence_worker_drain_tests=true
	;;
      --enable-turbulence-worker-failure-tests)
        enable_turbulence_worker_failure_tests=true
	;;
      --enable-turbulence-master-failure-tests)
        enable_turbulence_master_failure_tests=true
	;;
      --enable-turbulence-persistence-failure-tests)
        enable_turbulence_persistence_failure_tests=true
	;;
      --cidr-vars-file)
	cidr_vars_file="${value}"
	;;
      *)
        echo "$flag is not a valid flag"
        exit 1
        ;;
    esac
  done

  credhub_login $environment

  local director_name=$(bosh int "${VARFILE}" --path="/director_name")
  local iaas=$(bosh int "${VARFILE}" --path='/iaas')
  local routing_mode=$(bosh int "${VARFILE}" --path='/routing_mode')

  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' config <<-EOF
	{
	  "iaas": "$(bosh int ${VARFILE} --path=/iaas)",
	  "upgrade_tests": {
	    "include_multiaz": ${enable_multi_az_tests}
	  },
	  "bosh": {
	     "environment": "$(bosh int ${VARFILE} --path=/internal_ip)",
	     "ca_cert": $(bosh int ${CREDFILE} --path=/default_ca/ca --json | jq .Blocks[0]),
	     "client": "bosh_admin",
	     "client_secret": "$(bosh int ${CREDFILE} --path=/bosh_admin_client_secret)",
	     "deployment": "$deployment"
	  },
	  "turbulence": {
	    "host": "$(bosh int ${VARFILE} --path=/internal_ip)",
	    "port": 8080,
	    "username": "turbulence",
	    "password": "$(bosh int ${CREDFILE} --path=/turbulence_api_password 2>/dev/null)",
	    "ca_cert": $(bosh int ${CREDFILE} --path=/turbulence_api_ca/ca --json | jq .Blocks[0])
	  },
	  "turbulence_tests": {
	    "include_worker_drain": ${enable_turbulence_worker_drain_tests},
	    "include_worker_failure": ${enable_turbulence_worker_failure_tests},
	    "include_master_failure": ${enable_turbulence_master_failure_tests},
	    "include_persistence_failure": ${enable_turbulence_persistence_failure_tests},
	    "is_multiaz": ${enable_multi_az_tests}
	  },
	  "kubernetes": {
	    "master_host": "$(bosh int ${VARFILE} --path=/kubernetes_master_host)",
	    "master_port": $(bosh int ${VARFILE} --path=/kubernetes_master_port),
	    "tls_cert": $(bosh int <(credhub get -n "${director_name}/${deployment}/tls-kubernetes" --output-json) --path='/value/certificate' --json | jq .Blocks[0]),
	    "tls_private_key": $(bosh int <(credhub get -n "${director_name}/${deployment}/tls-kubernetes" --output-json) --path='/value/private_key' --json | jq .Blocks[0]),
	    "cluster_ip_range": "$(bosh int ${ROOT}/${cidr_vars_file} --path=/service_cluster_cidr)",
	    "kubernetes_service_ip": "$(bosh int ${ROOT}/${cidr_vars_file} --path=/first_ip_of_service_cluster_cidr)",
	    "kube_dns_ip": "$(bosh int ${ROOT}/${cidr_vars_file} --path=/kubedns_service_ip)",
	    "pod_ip_range": "$(bosh int ${ROOT}/${cidr_vars_file} --path=/pod_network_cidr)"
	  },
	  "timeout_scale": $(bosh int ${VARFILE} --path=/timeout_scale 2>/dev/null || echo 1),
	  "cfcr": {
	    "deployment_path": "${ROOT}/git-kubo-deployment"
	  }
	}
	EOF
  set -e

  if [[ "${iaas}" == "aws" ]]; then
    config="$(echo ${config} | jq \
      --arg access_key_id "$(bosh int "${VARFILE}" --path=/access_key_id)" \
      --arg secret_access_key "$(bosh int "${VARFILE}" --path=/secret_access_key)" \
      --arg region "$(bosh int "${VARFILE}" --path=/region)" \
      --arg ingress_group_id "$(bosh int "${VARFILE}" --path=/default_security_groups/0)" \
      '. +
      {
	"aws": {
	  "access_key_id": $access_key_id,
	  "secret_access_key": $secret_access_key,
	  "region": $region,
	  "ingress_group_id": $ingress_group_id
	}
      }')"
  fi

  echo "$config" | jq .
}

main() {
  generate_test_config $@
}

main $@
