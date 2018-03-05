#!/bin/bash

set -euxo pipefail

stemcell_version=$(cat ./gcp-stemcell/version)
./git-kubo-deployment/bin/update_stemcell "${stemcell_version}"

cp -a git-kubo-deployment/. git-kubo-deployment-with-updated-stemcell

pushd git-kubo-deployment-with-updated-stemcell

  git config user.email "ci-bot@localhost"
  git config user.name "CI Bot"

  git checkout master
  git add .
  git diff-index --quiet HEAD || git commit -m "Update stemcell version to $stemcell_version"

popd
