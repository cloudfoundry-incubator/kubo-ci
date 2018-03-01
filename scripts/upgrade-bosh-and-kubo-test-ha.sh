#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/utils.sh"

KUBO_ENVIRONMENT_DIR=$1
DEPLOYMENT_NAME=$2

tmpfile=$(mktemp)
$DIR/generate-test-config.sh "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}" > "${tmpfile}"
export CONFIG="${tmpfile}"

IAAS=$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path='/iaas')

if [[ "${IAAS}" == "aws" ]]; then
  aws configure set aws_access_key_id "$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/access_key_id)"
  aws configure set aws_secret_access_key  "$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/secret_access_key)"
  aws configure set default.region "$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/region)"
  AWS_INGRESS_GROUP_ID=$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/default_security_groups/0)
  export AWS_INGRESS_GROUP_ID
fi

ginkgo -progress -v -failFast "$DIR/../src/tests/upgrade-tests"
