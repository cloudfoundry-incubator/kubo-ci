#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

BASE_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)

verify_args() {
  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' usage <<-EOF
	Usage: $(basename "$0") [-h] environment deployment-name [--skip-addons-tests] [--skip-cloud-agnostic-tests]

	Help Options:
          -h                             show this help text
          --skip-addons-tests            skip the add ons tests
          --skip-cloud-agnostic-tests   skip the cloud agnostic tests
	EOF
  set -e

  while getopts ':h' option; do
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

execute_cloud_agnostic_tests() {
  local routing_mode="$1"
  local skip_addons_tests="$2"
  local cloud_agnostic_tests=("pod_logs" "generic" "oss_only" "api_extensions")
  local ginkgo_flags

  if ! [[ -z "$skip_addons_tests" ]]; then
    ginkgo_flags="--skip=check\ apply-specs"
  fi

  if [[ ${routing_mode} == "cf" ]]; then
    ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/cloudfoundry"
  fi

  for test in "${cloud_agnostic_tests[@]}"; do
    ginkgo -progress -v "$ginkgo_flags" "$BASE_DIR/src/tests/integration-tests/${test}"
  done

}

execute_cloud_specific_tests(){
  local routing_mode="$1"
  local iaas="$2"

  if [[ ${routing_mode} == "iaas" ]]; then
    case "${iaas}" in
      aws)
        aws configure set aws_access_key_id "$(bosh-cli int "${environment}/director.yml" --path=/access_key_id)"
        aws configure set aws_secret_access_key  "$(bosh-cli int "${environment}/director.yml" --path=/secret_access_key)"
        aws configure set default.region "$(bosh-cli int "${environment}/director.yml" --path=/region)"
        AWS_INGRESS_GROUP_ID=$(bosh-cli int "${environment}/director.yml" --path=/default_security_groups/0)
        export AWS_INGRESS_GROUP_ID
        ;;
    esac

    ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/workload/k8s_lbs"
  fi

  if [[ "${iaas}" != "openstack" ]]; then
      ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/persistent_volume"
  fi
}

run_tests() {
  local environment="$1"
  local deployment="$2"
  local skip_addons_tests=""
  local skip_cloud_agnostic_tests=""

  shift 2
  for flag in "$@"; do
    case "$flag" in
      --skip-addons-tests)
        skip_addons_tests=true
        ;;
      --skip-cloud-agnostic-tests)
        skip_cloud_agnostic_tests=true
        ;;
      *)
        echo "$flag is not a valid flag"
        exit 1
        ;;
    esac
  done

  local iaas=$(bosh-cli int "$environment/director.yml" --path='/iaas')
  local routing_mode=$(bosh-cli int "$environment/director.yml" --path='/routing_mode')

  local tmpfile=$(mktemp)
  $BASE_DIR/scripts/generate-test-config.sh $environment $deployment > $tmpfile
  export CONFIG=$tmpfile

  if [[ -z "$skip_cloud_agnostic_tests" ]]; then
    execute_cloud_agnostic_tests "${routing_mode}" "${skip_addons_tests}"
  fi

  execute_cloud_specific_tests "${routing_mode}" "${iaas}"

  return 0
}

main() {
  verify_args "$@"
  run_tests "$@"
}

main "$@"
