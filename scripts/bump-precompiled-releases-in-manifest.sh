#!/usr/bin/env bash

set -eu

cp -r git-kubo-deployment/. git-kubo-deployment-output

stemcell_os=$(bosh int git-kubo-deployment/manifests/cfcr.yml --path=/stemcells/0/os)
stemcell_version=$(bosh int git-kubo-deployment/manifests/cfcr.yml --path=/stemcells/0/version)

echo > bump-precompiled-release.yml

for release in $RELEASE_LIST
do
  if [[ "$release" == "kubo" ]]; then
    version="$(cat kubo-version/version)"
  else
    version=$(bosh int git-kubo-deployment/manifests/cfcr.yml -o git-kubo-deployment/manifests/ops-files/non-precompiled-releases.yml "--path=/releases/name=$release/version")
  fi
  release_path=$(ls compiled-releases/$release-*.tgz)
  sha1=$(sha1sum ${release_path} | awk '{print $1}')
  url="https://storage.googleapis.com/kubo-precompiled-releases/$(basename ${release_path})"

cat >> bump-precompiled-releases.yml <<EOF
- type: replace
  path: /releases/name=$release
  value:
    name: $release
    version: $version
    sha1: $sha1
    url: $url
    stemcell:
      os: $stemcell_os
      version: $stemcell_version
EOF

done

bosh int git-kubo-deployment/manifests/cfcr.yml \
  -o bump-precompiled-releases.yml > git-kubo-deployment-output/manifests/cfcr.yml

pushd git-kubo-deployment-output
git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"

if [ -n "$(git status --porcelain)" ]; then
    git add manifests/cfcr.yml
    git commit -m "Precompile $RELEASE_LIST against $stemcell_os/$stemcell_version"
else
    echo "No changes to commit."
fi


