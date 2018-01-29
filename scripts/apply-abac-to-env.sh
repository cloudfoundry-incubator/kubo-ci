#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"


cp "$PWD/kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"


cp -r kubo-lock/* kubo-lock-with-abac/


sed -iE 's/rbac$/abac/g' kubo-lock-with-abac/metadata
