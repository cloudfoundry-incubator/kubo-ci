#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
setup_env "${KUBO_ENVIRONMENT_DIR}"

echo "Succeed"
