#!/usr/bin/env bash

set -eu

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" && pwd)
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

mkdir -p "${KUBO_ENVIRONMENT_DIR}"
cp "${ROOT}/gaffer-director-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "${ROOT}/gaffer-director-metadata/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

BOSH_ENV="${KUBO_ENVIRONMENT_DIR}" source "${ROOT}/git-kubo-ci/scripts/set_bosh_environment"

STEMCELL_OS=$(bosh int ${ROOT}/git-kubo-deployment/manifests/cfcr.yml --path=/stemcells/0/os)
STEMCELL_VERSION=$(bosh int ${ROOT}/git-kubo-deployment/manifests/cfcr.yml --path=/stemcells/0/version)

for release in $RELEASE_LIST
do
  if [[ "$release" == "kubo" ]]; then
    RELEASE_VERSION="$(cat kubo-version/version)"
  else
    RELEASE_VERSION="$(bosh int ${ROOT}/git-kubo-deployment/manifests/cfcr.yml --path=/releases/name=${release}/version)"
  fi
  bosh -d compilation export-release "$release/$RELEASE_VERSION" "$STEMCELL_OS/$STEMCELL_VERSION"

  mv *.tgz "compiled-releases/$release-$RELEASE_VERSION-$STEMCELL_OS-$STEMCELL_VERSION.tgz"
done
