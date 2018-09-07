#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

export DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

export GOPATH="${ROOT}/git-kubo-ci"

main() {
  local tmpfile release_tarball

  source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
  setup_env "${KUBO_ENVIRONMENT_DIR}"

  BOSH_ENV="${KUBO_ENVIRONMENT_DIR}" source "${ROOT}/git-kubo-deployment/bin/set_bosh_environment"

  release_tarball=$(find "${ROOT}/gcs-kubo-release-tarball/" -name "*kubo-*.tgz" | head -n1)

  bosh upload-release "$release_tarball"
  bosh upload-stemcell "${ROOT}/stemcell/stemcell.tgz"

  if [[ "$INTERNET_ACCESS" != "false" ]]; then
    tmpfile="$(mktemp)"
    echo "CONFIG=${tmpfile}"

    "${ROOT}/git-kubo-ci/scripts/generate-test-config.sh" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}" > "${tmpfile}"
    BOSH_DEPLOY_COMMAND="$ROOT/bosh-command/bosh-deploy.sh" CONFIG="${tmpfile}" ginkgo -r -v -progress "${ROOT}/git-kubo-ci/src/tests/upgrade-tests/"
  else
    $ROOT/bosh-command/bosh-deploy.sh
  fi
}

main
