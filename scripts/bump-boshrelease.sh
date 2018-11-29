#!/bin/bash

set -exu -o pipefail

cp -r git-kubo-deployment/. git-kubo-deployment-output

release="${RELEASE_NAME}"
version="$(cat boshrelease/version)"

if [ -z ${REPO_URL} ]; then
 url="$(cat boshrelease/url)"
else
 url="${REPO_URL}/releases/download/v${version}/${release}-release-${version}.tgz"
fi

if [ -f "boshrelease/sha1" ]; then
  sha1="$(cat boshrelease/sha1)"
else
  sha1=$(sha1sum boshrelease/*.tgz | awk '{print $1}')
fi

cat > update-$release-release.yml <<EOF
- type: replace
  path: /0/value/name=$release
  value:
    name: $release
    version: $version
    sha1: $sha1
    url: $url
EOF

bosh int git-kubo-deployment/manifests/ops-files/non-precompiled-releases.yml \
  -o update-$release-release.yml > git-kubo-deployment-output/manifests/ops-files/non-precompiled-releases.yml

pushd git-kubo-deployment-output
git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"

if [ -n "$(git status --porcelain)" ]; then
    git add manifests/ops-files/non-precompiled-releases.yml
    git commit -m "Bump $release to version $version in non-precompiled-releases ops-file"
else
    echo "No changes to commit."
fi
