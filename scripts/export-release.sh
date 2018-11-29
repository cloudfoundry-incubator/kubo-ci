#!/usr/bin/env bash

set -eu

source git-kubo-ci/scripts/set-bosh-env gcs-source-json/source.json

STEMCELL_OS=$(bosh int git-kubo-deployment/manifests/cfcr.yml --path=/stemcells/0/os)
STEMCELL_VERSION=$(bosh int git-kubo-deployment/manifests/cfcr.yml --path=/stemcells/0/version)

for release in $RELEASE_LIST
do
  if [[ "$release" == "kubo" ]]; then
    RELEASE_VERSION="$(cat kubo-version/version)"
  else
    RELEASE_VERSION="$(bosh int git-kubo-deployment/manifests/cfcr.yml -o git-kubo-deployment/manifests/ops-files/non-precompiled-releases.yml --path=/releases/name=${release}/version)"
  fi
  bosh -d compilation export-release "$release/$RELEASE_VERSION" "$STEMCELL_OS/$STEMCELL_VERSION"

  mv *.tgz "compiled-releases/"
done
