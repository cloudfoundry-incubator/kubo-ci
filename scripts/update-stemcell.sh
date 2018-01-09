#!/bin/bash

set -euxo pipefail

stemcell_version=$(cat ./gcp-stemcell/version)
./git-kubo-deployment/bin/update_stemcell "${stemcell_version}"

pushd git-kubo-deployment

  git config user.email "ci-bot@localhost"
  git config user.name "CI Bot"

  git commit -am "Update stemcell version to $stemcell_version"

popd
