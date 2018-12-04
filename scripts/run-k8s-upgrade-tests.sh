#!/usr/bin/env bash

set -eu -o pipefail

export ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
export KUBO_ENVIRONMENT_DIR="${ROOT}/environment"
export GOPATH="${ROOT}/git-kubo-ci"

main() {
  local release_tarball

  source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
  setup_env "${KUBO_ENVIRONMENT_DIR}"

  BOSH_ENV="${KUBO_ENVIRONMENT_DIR}" source "${ROOT}/git-kubo-ci/scripts/set_bosh_environment"

  release_tarball=$(find "${ROOT}/gcs-kubo-release-tarball/" -name "*kubo-*.tgz" | head -n1)
  bosh upload-release "$release_tarball"

  if [[ "$INTERNET_ACCESS" != "false" ]]; then
    BOSH_DEPLOY_COMMAND="$ROOT/bosh-command/bosh-deploy.sh" ginkgo -r -v -progress "${ROOT}/git-kubo-ci/src/tests/upgrade-tests/"
  else
    $ROOT/bosh-command/bosh-deploy.sh
  fi
}

main
