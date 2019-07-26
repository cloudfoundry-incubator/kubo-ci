#!/usr/bin/env bash

set -eu

stemcell_alias=${stemcell_alias:-default}
source git-kubo-ci/scripts/set-bosh-env gcs-source-json/source.json

pushd git-kubo-deployment/manifests
STEMCELL_OS=$(bosh int cfcr.yml --ops-file ops-files/windows/add-worker.yml --path=/stemcells/alias=${stemcell_alias}/os)
STEMCELL_VERSION=$(bosh int cfcr.yml --ops-file ops-files/windows/add-worker.yml --path=/stemcells/alias=${stemcell_alias}/version )
popd

jobs=""
for job in $JOBS_LIST
do
  jobs="${jobs} --job $job"
done

for release in $RELEASE_LIST
do
  if [[ $release == kubo* ]]; then
    RELEASE_VERSION="$(cat kubo-version/version)"
  else
    RELEASE_VERSION="$(bosh int git-kubo-deployment/manifests/cfcr.yml -o git-kubo-deployment/manifests/ops-files/non-precompiled-releases.yml --path=/releases/name=${release}/version)"
  fi
  bosh -d compilation-${stemcell_alias} export-release "$release/$RELEASE_VERSION" "$STEMCELL_OS/$STEMCELL_VERSION" ${jobs}

  mv *.tgz "compiled-releases/"
done
