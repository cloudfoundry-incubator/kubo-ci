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
  sha1=$(sha1sum compiled-releases/$release*.tgz | awk '{print $1}')
  url="https://storage.cloud.google.com/kubo-precompiled-releases/$release-$version-$stemcell_os-$stemcell_version.tgz"
cat >> bump-precompiled-releases.yml <<EOF
- type: replace
  path: /releases/name=$release
  value:
    name: $release
    version: $version
    sha1: $sha1
    url: $url
EOF

done

bosh int git-kubo-deployment/manifests/cfcr.yml \
  -o bump-precompiled-releases.yml > git-kubo-deployment-output/manifests/cfcr.yml

pushd git-kubo-deployment-output
git config --global user.email "cfcr+cibot@pivotal.io"
git config --global user.name "CFCR CI BOT"

if git diff-index --quiet HEAD; then
    echo "No changes to commit."
else
    git add manifests/cfcr.yml
    git commit -m "Bump $release to version $version"
fi


