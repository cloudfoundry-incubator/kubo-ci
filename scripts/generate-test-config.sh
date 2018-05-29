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
		--enable-addons-tests                         [env:ENABLE_ADDONS_TESTS]
		--enable-api-extensions-tests                 [env:ENABLE_API_EXTENSIONS_TESTS]
		--enable-generic-tests                        [env:ENABLE_GENERIC_TESTS]
		--enable-multi-az-tests                       [env:ENABLE_MULTI_AZ_TESTS]
		--enable-oss-only-tests                       [env:ENABLE_OSS_ONLY_TESTS]
		--enable-persistent-volume-tests              [env:ENABLE_PERSISTENT_VOLUME_TESTS]
		--enable-certificate-tests                    [env:ENABLE_CERTIFICATE_TESTS]

		--conformance_release_version=<some-value>    [env:CONFORMANCE_RELEASE_VERSION]
		--conformance_results_dir=<some-value>        [env:CONFORMANCE_RESULTS_DIR]

		--new-bosh-stemcell-version=<some-value>      [env:NEW_BOSH_STEMCELL_VERSION]

		--enable-turbulence-master-failure-tests      [env:ENABLE_TURBULENCE_MASTER_FAILURE_TESTS]
		--enable-turbulence-persistence-failure-tests [env:ENABLE_TURBULENCE_PERSISTENCE_FAILURE_TESTS]
		--enable-turbulence-worker-drain-tests        [env:ENABLE_TURBULENCE_WORKER_DRAIN_TESTS]
		--enable-turbulence-worker-failure-tests      [env:ENABLE_TURBULENCE_WORKER_FAILURE_TESTS]
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
  local enable_addons_tests="${ENABLE_ADDONS_TESTS:-false}"
  local enable_api_extensions_tests="${ENABLE_API_EXTENSIONS_TESTS:-false}"
  local enable_generic_tests="${ENABLE_GENERIC_TESTS:-false}"
  local enable_certificate_tests="${ENABLE_CERTIFICATE_TESTS:-false}"
  local enable_multi_az_tests="${ENABLE_MULTI_AZ_TESTS:-false}"
  local enable_oss_only_tests="${ENABLE_OSS_ONLY_TESTS:-false}"
  local enable_persistent_volume_tests="${ENABLE_PERSISTENT_VOLUME_TESTS:-false}"
  local conformance_release_version="${CONFORMANCE_RELEASE_VERSION:-dev}"
  local conformance_results_dir="${CONFORMANCE_RESULTS_DIR:-/tmp}"
  local new_bosh_stemcell_version="${NEW_BOSH_STEMCELL_VERSION:-""}"
  local enable_turbulence_worker_drain_tests="${ENABLE_TURBULENCE_WORKER_DRAIN_TESTS:-false}"
  local enable_turbulence_worker_failure_tests="${ENABLE_TURBULENCE_WORKER_FAILURE_TESTS:-false}"
  local enable_turbulence_master_failure_tests="${ENABLE_TURBULENCE_MASTER_FAILURE_TESTS:-false}"
  local enable_turbulence_persistence_failure_tests="${ENABLE_TURBULENCE_PERSISTENCE_FAILURE_TESTS:-false}"

  shift 2
  for arg in "$@"; do
    local flag="${arg%%=*}"
    local value="${arg##*=}"
    case "$flag" in
      --enable-addons-tests)
	enable_addons_tests=true
	;;
      --enable-api-extensions-tests)
	enable_api_extensions_tests=true
	;;
      --enable-generic-tests)
	enable_generic_tests=true
	;;
      --enable-certificate-tests)
	enable_certificate_tests=true
	;;
      --enable-multi-az-tests)
	enable_multi_az_tests=true
	;;
      --enable-oss-only-tests)
	enable_oss_only_tests=true
	;;
      --enable-persistent-volume-tests)
	enable_persistent_volume_tests=true
	;;
      --conformance_release_version)
	conformance_release_version="${value}"
	;;
      --conformance_results_dir)
	conformance_results_dir="${value}"
	;;
      --new-bosh-stemcell-version)
	new_bosh_stemcell_version="${value}"
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
  local iaas=$(bosh int "$environment/director.yml" --path='/iaas')
  local routing_mode=$(bosh int "$environment/director.yml" --path='/routing_mode')

  local enable_rbac_tests="false"
  if [[ ${authorization_mode} == "rbac" ]]; then
    enable_rbac_tests="true"
  fi

  local enable_cloudfoundry_tests="false"
  if [[ ${routing_mode} == "cf" ]]; then
    enable_cloudfoundry_tests="true"
  fi

  local enable_iaas_k8s_lb_tests="false"
  if [[ ${routing_mode} == "iaas" ]]; then
    enable_iaas_k8s_lb_tests="true"
  fi

  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' config <<-EOF
	{
	  "iaas": "$(bosh int $director_yml --path=/iaas)",
	  "integration_tests": {
	    "include_certificates": ${enable_certificate_tests},
	    "include_cloudfoundry": ${enable_cloudfoundry_tests},
	    "include_generic": ${enable_generic_tests},
	    "include_k8s_lb": ${enable_iaas_k8s_lb_tests},
	    "include_multiaz": ${enable_multi_az_tests},
	    "include_oss_only": ${enable_oss_only_tests},
	    "include_persistent_volume": ${enable_persistent_volume_tests},
	    "include_rbac": ${enable_rbac_tests}
	  },
	  "upgrade_tests": {
	    "include_multiaz": ${enable_multi_az_tests}
	  },
	  "conformance": {
	    "results_dir": "${conformance_results_dir}",
	    "release_version": "${conformance_release_version}"
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
