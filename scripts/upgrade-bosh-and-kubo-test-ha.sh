#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"
. "$DIR/lib/utils.sh"

KUBO_ENVIRONMENT_DIR=$1
DEPLOYMENT_NAME=$2

tmpfile=$(mktemp)
$DIR/generate-test-config.sh "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}" > "${tmpfile}"
export CONFIG="${tmpfile}"

ginkgo -progress -v "$DIR/../src/tests/upgrade-tests"
