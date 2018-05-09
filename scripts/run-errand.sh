#!/bin/bash

set -euo pipefail

KUBO_CI="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../../" && pwd )"
source "$KUBO_CI/scripts/bosh_director_creds.sh"

bosh -n -d "${DEPLOYMENT_NAME}" run-errand --instance "${INSTANCE}" "${ERRAND_NAME}"
