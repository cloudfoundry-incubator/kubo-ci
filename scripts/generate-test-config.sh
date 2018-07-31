#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"

verify_args() {
  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' usage <<-EOF
	Usage: $(basename "$0") [-h] environment deployment-name

	Options:
		-h                                            show this help text
		--enable-multi-az-tests                       [env:ENABLE_MULTI_AZ_TESTS]

		--enable-turbulence-master-failure-tests      [env:ENABLE_TURBULENCE_MASTER_FAILURE_TESTS]
		--enable-turbulence-persistence-failure-tests [env:ENABLE_TURBULENCE_PERSISTENCE_FAILURE_TESTS]
		--enable-turbulence-worker-drain-tests        [env:ENABLE_TURBULENCE_WORKER_DRAIN_TESTS]
		--enable-turbulence-worker-failure-tests      [env:ENABLE_TURBULENCE_WORKER_FAILURE_TESTS]
		--cidr-vars-file=<some-path>                  [env:CIDR_VARS_FILE]
	EOF
  set -e

  while getopts ':h:' option; do
    case "$option" in
      h) echo "$usage"
         exit 0
         ;;
     \?) printf "Illegal option: -%s\n" "$OPTARG" >&2
         echo "$usage" >&2
         exit 64
         ;;
    esac
  done
  shift $((OPTIND - 1))
  if [[ $# -lt 2 ]]; then
    echo "$usage" >&2
    exit 64
  fi
}

credhub_login() {
  local environment="$1"

  local credhub_admin_secret=$(bosh int "${environment}/creds.yml" --path="/credhub_admin_client_secret")
  local credhub_api_url="https://$(bosh int "${environment}/director.yml" --path="/internal_ip"):8844"

  credhub login --client-name credhub-admin --client-secret "${credhub_admin_secret}" \
    -s "${credhub_api_url}" \
    --ca-cert <(bosh int "${environment}/creds.yml" --path="/credhub_tls/ca") \
    --ca-cert <(bosh int "${environment}/creds.yml" --path="/uaa_ssl/ca") 1>/dev/null
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

  local director_yml="$environment"/director.yml
  local creds_yml="$environment"/creds.yml

  credhub_login $environment

  local director_name=$(bosh int "${environment}/director.yml" --path="/director_name")
  local iaas=$(bosh int "$environment/director.yml" --path='/iaas')
  local routing_mode=$(bosh int "$environment/director.yml" --path='/routing_mode')

  local enable_cloudfoundry_tests="false"
  if [[ ${routing_mode} == "cf" ]]; then
    enable_cloudfoundry_tests="true"
  fi

  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' config <<-EOF
	{
	  "iaas": "$(bosh int $director_yml --path=/iaas)",
	  "upgrade_tests": {
	    "include_multiaz": ${enable_multi_az_tests}
	  },
	  "bosh": {
	     "environment": "$(bosh int $director_yml --path=/internal_ip)",
	     "ca_cert": $(bosh int $creds_yml --path=/default_ca/ca --json | jq .Blocks[0]),
	     "client": "bosh_admin",
	     "client_secret": "$(bosh int $creds_yml --path=/bosh_admin_client_secret)",
	     "deployment": "$deployment"
	  },
	  "turbulence": {
	    "host": "$(bosh int $director_yml --path=/internal_ip)",
	    "port": 8080,
	    "username": "turbulence",
	    "password": "$(bosh int $creds_yml --path=/turbulence_api_password 2>/dev/null)",
	    "ca_cert": $(bosh int $creds_yml --path=/turbulence_api_ca/ca --json | jq .Blocks[0])
	  },
	  "turbulence_tests": {
	    "include_worker_drain": ${enable_turbulence_worker_drain_tests},
	    "include_worker_failure": ${enable_turbulence_worker_failure_tests},
	    "include_master_failure": ${enable_turbulence_master_failure_tests},
	    "include_persistence_failure": ${enable_turbulence_persistence_failure_tests},
	    "is_multiaz": ${enable_multi_az_tests}
	  },
	  "kubernetes": {
	    "master_host": "$(bosh int $director_yml --path=/kubernetes_master_host)",
	    "master_port": $(bosh int $director_yml --path=/kubernetes_master_port),
	    "tls_cert": $(bosh int <(credhub get -n "${director_name}/${deployment}/tls-kubernetes" --output-json) --path='/value/certificate' --json | jq .Blocks[0]),
	    "tls_private_key": $(bosh int <(credhub get -n "${director_name}/${deployment}/tls-kubernetes" --output-json) --path='/value/private_key' --json | jq .Blocks[0]),
	    "cluster_ip_range": "$(bosh int ${ROOT}/${cidr_vars_file} --path=/service_cluster_cidr)",
	    "kubernetes_service_ip": "$(bosh int ${ROOT}/${cidr_vars_file} --path=/first_ip_of_service_cluster_cidr)",
	    "kube_dns_ip": "$(bosh int ${ROOT}/${cidr_vars_file} --path=/kubedns_service_ip)",
	    "pod_ip_range": "$(bosh int ${ROOT}/${cidr_vars_file} --path=/pod_network_cidr)"
	  },
	  "timeout_scale": $(bosh int $director_yml --path=/timeout_scale 2>/dev/null || echo 1),
	  "cfcr": {
	    "deployment_path": "${ROOT}/git-kubo-deployment"
	  }
	}
	EOF
  set -e

  if [[ "${routing_mode}" == "iaas" && "${iaas}" == "aws" ]]; then
    config="$(echo ${config} | jq \
      --arg access_key_id "$(bosh int "${environment}/director.yml" --path=/access_key_id)" \
      --arg secret_access_key "$(bosh int "${environment}/director.yml" --path=/secret_access_key)" \
      --arg region "$(bosh int "${environment}/director.yml" --path=/region)" \
      --arg ingress_group_id "$(bosh int "${environment}/director.yml" --path=/default_security_groups/0)" \
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
  verify_args $@
  generate_test_config $@
}

main $@
