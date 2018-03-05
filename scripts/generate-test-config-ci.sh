#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"

verify_args() {
  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' usage <<-EOF
	Usage: $(basename "$0") [-h] environment deployment-name --enable-multi-az-tests

	Help Options:
		-h  show this help text
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
  local enable_addons_tests="false"
  local enable_multi_az_tests="false"

  shift 2
  for flag in "$@"; do
    case "$flag" in
      --enable-addons-tests)
        enable_addons_tests=true
        ;;
      --enable-multi-az-tests)
        enable_multi_az_tests=true
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
  local authorization_mode=$(bosh int "${environment}/director.yml" --path='/authorization_mode')
  local routing_mode=$(bosh int "$environment/director.yml" --path='/routing_mode')

  local enable_rbac_tests="false"
  if [[ ${authorization_mode} == "rbac" ]]; then
    enable_rbac_tests="true"
  fi

  local enable_cloudfoundry_tests="false"
  if [[ ${routing_mode} == "cf" ]]; then
    enable_cloudfoundry_tests="true"
  fi

  local new_bosh_stemcell_version=""
  if [[ -f "${ROOT}/new-bosh-stemcell/version" ]]; then
    new_bosh_stemcell_version="$(cat ${ROOT}/new-bosh-stemcell/version)"
  fi

  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' config <<-EOF
	{
	  "test_suites": {
	    "include_api_extensions": true,
	    "include_generic": true,
	    "include_addons": ${enable_addons_tests},
	    "include_oss_only": true,
	    "include_pod_logs": true,
	    "include_rbac": ${enable_rbac_tests},
	    "include_cloudfoundry": ${enable_cloudfoundry_tests},
	    "include_multiaz": ${enable_multi_az_tests},
	    "include_workload": true,
	    "include_persistent_volume": true
	  },
	  "bosh": {
	     "iaas": "$(bosh int $director_yml --path=/iaas)",
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
	  "cf": {
	    "apps_domain": "$(bosh int $director_yml --path=/routing_cf_app_domain_name 2>/dev/null)"
	  },
	  "kubernetes": {
	    "authorization_mode": "$(bosh int $director_yml --path=/authorization_mode)",
	    "master_host": "$(bosh int $director_yml --path=/kubernetes_master_host)",
	    "master_port": $(bosh int $director_yml --path=/kubernetes_master_port),
	    "path_to_kube_config": "$HOME/.kube/config",
	    "tls_cert": $(bosh int <(credhub get -n "${director_name}/${deployment}/tls-kubernetes" --output-json) --path='/value/certificate' --json | jq .Blocks[0]),
	    "tls_private_key": $(bosh int <(credhub get -n "${director_name}/${deployment}/tls-kubernetes" --output-json) --path='/value/private_key' --json | jq .Blocks[0])
	  },
	  "timeout_scale": $(bosh int $director_yml --path=/timeout_scale 2>/dev/null || echo 1),
	  "cfcr": {
	    "deployment_path": "${ROOT}/git-kubo-deployment",
	    "upgrade_to_stemcell_version": "${new_bosh_stemcell_version}"
	  }
	}
	EOF
  set -e

  echo "$config"
}

main() {
  verify_args $@
  generate_test_config $@
}

main $@
