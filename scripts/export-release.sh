#!/usr/bin/env bash

set -eu

stemcell_alias=${stemcell_alias:-default}
source git-kubo-ci/scripts/set-bosh-env gcs-source-json/source.json

pushd git-kubo-deployment/manifests
STEMCELL_OS=$(bosh int cfcr.yml --ops-file ops-files/windows/add-worker.yml --path=/stemcells/alias=${stemcell_alias}/os)
STEMCELL_VERSION=$(bosh stemcells --json | jq .Tables[0].Rows | jq -r "map(select(.os | contains(\"${STEMCELL_OS}\"))) | max_by(.version) | .version | sub(\"\\\*\"; \"\")")
popd

for release in $RELEASE_LIST
do
  if [[ $release == kubo* ]]; then
    RELEASE_VERSION="$(cat kubo-version/version)"
  elif [[ "$release" == "kubernetes*" ]]; then
    RELEASE_VERSION="$(cat kubernetes-version/version)"
  else
    RELEASE_VERSION="$(bosh int git-kubo-deployment/manifests/cfcr.yml -o git-kubo-deployment/manifests/ops-files/non-precompiled-releases.yml --path=/releases/name=${release}/version)"
  fi
  bosh -d compilation-${stemcell_alias} export-release "$release/$RELEASE_VERSION" "$STEMCELL_OS/$STEMCELL_VERSION"

  mv *.tgz "compiled-releases/"
done
