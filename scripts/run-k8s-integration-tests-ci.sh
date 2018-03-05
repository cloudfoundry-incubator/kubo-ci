#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

BASE_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)

ENVIRONMENT=${ENVIRONMENT:-"default"}
DEPLOYMENT=${DEPLOYMENT:-"default"}
ENABLE_ADDONS_TESTS=${ENABLE_ADDONS_TESTS}
ENABLE_MULTI_AZ_TESTS=${ENABLE_MULTI_AZ_TESTS}
IAAS=${IAAS}

execute_cloud_specific_tests(){
  local routing_mode="$1"
  local iaas="$2"

  if [[ ${routing_mode} == "iaas" ]]; then
    case "${iaas}" in
      aws)
        aws configure set aws_access_key_id "$(bosh int "${environment}/director.yml" --path=/access_key_id)"
        aws configure set aws_secret_access_key  "$(bosh int "${environment}/director.yml" --path=/secret_access_key)"
        aws configure set default.region "$(bosh int "${environment}/director.yml" --path=/region)"
        AWS_INGRESS_GROUP_ID=$(bosh int "${environment}/director.yml" --path=/default_security_groups/0)
        export AWS_INGRESS_GROUP_ID
        ;;
    esac

    ginkgo -r -progress -v "$BASE_DIR/src/tests/integration-tests/workload"
  fi

  ginkgo -progress -v "$BASE_DIR/src/tests/integration-tests/persistent_volume"
}

run_tests() {
  local tmpfile=$(mktemp)
  $BASE_DIR/scripts/generate-test-config.sh $environment $deployment > $tmpfile
  export CONFIG=$tmpfile

  # local iaas=$(bosh int "$environment/director.yml" --path='/iaas')
  # local routing_mode=$(bosh int "$environment/director.yml" --path='/routing_mode')
  # execute_cloud_specific_tests "${routing_mode}" "${iaas}"

  ginkgo -r -progress -dryRun  -v "$BASE_DIR/src/tests/integration-tests/"
  return 0
}

run_tests "$@"
