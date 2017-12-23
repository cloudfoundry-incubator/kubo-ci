#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu

verify_args() {
  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' usage <<-EOF
	Usage: $(basename "$0") [-h] environment deployment-name
	
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

  local credhub_user_password=$(bosh-cli int "${environment}/creds.yml" --path="/credhub_cli_password")
  local credhub_api_url="https://$(bosh-cli int "${environment}/director.yml" --path="/internal_ip"):8844"

  credhub login -u credhub-cli -p "${credhub_user_password}" \
    -s "${credhub_api_url}" \
    --ca-cert <(bosh-cli int "${environment}/creds.yml" --path="/credhub_tls/ca") \
    --ca-cert <(bosh-cli int "${environment}/creds.yml" --path="/uaa_ssl/ca") 1>/dev/null
}

generate_test_config() {
  local environment="$1"
  local deployment="$2"

  local director_yml="$environment"/director.yml
  local creds_yml="$environment"/creds.yml

  credhub_login $environment

  local director_name=$(bosh-cli int "${environment}/director.yml" --path="/director_name")

  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' config <<-EOF
	{
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
	    "authorization_mode": "$(bosh-cli int $director_yml --path=/authorization_mode)",
	    "master_host": "$(bosh int $director_yml --path=/kubernetes_master_host)",
	    "master_port": $(bosh int $director_yml --path=/kubernetes_master_port),
	    "path_to_kube_config": "$HOME/.kube/config",
	    "tls_cert": $(bosh-cli int <(credhub get -n "${director_name}/${deployment}/tls-kubernetes" --output-json) --path='/value/certificate' --json | jq .Blocks[0]),
	    "tls_private_key": $(bosh-cli int <(credhub get -n "${director_name}/${deployment}/tls-kubernetes" --output-json) --path='/value/private_key' --json | jq .Blocks[0])
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
