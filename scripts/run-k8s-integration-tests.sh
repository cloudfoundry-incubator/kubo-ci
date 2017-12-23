#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

BASE_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)

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

run_tests() {
  local environment="$1"
  local deployment="$2"

  local iaas=$(bosh-cli int "$environment/director.yml" --path='/iaas')
  local routing_mode=$(bosh-cli int "$environment/director.yml" --path='/routing_mode')
  local director_name=$(bosh-cli int "$environment/director.yml" --path='/director_name')

  local tmpfile=$(mktemp)
  $BASE_DIR/scripts/generate-test-config.sh $environment $deployment > $tmpfile
  export CONFIG=$tmpfile

  if [[ ${routing_mode} == "cf" ]]; then
    ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/cloudfoundry"
  elif [[ ${routing_mode} == "iaas" ]]; then
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

  ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/pod_logs"
  ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/generic"
  ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/oss_only"
  ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/api_extensions"

  if [[ "${iaas}" != "openstack" ]]; then
      ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/persistent_volume"
  fi

  return 0
}

main() {
  verify_args "$@"
  run_tests "$@"
}

main "$@"
