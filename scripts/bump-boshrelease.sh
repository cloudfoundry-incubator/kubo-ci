#!/bin/bash

set -exu -o pipefail

cp -r git-kubo-deployment/. git-kubo-deployment-output

release="${RELEASE_NAME}"
version="$(cat boshrelease/version)"
url="$(cat boshrelease/url)"
sha1="$(cat boshrelease/sha1)"

cat > update-$release-release.yml <<EOF
- type: replace
  path: /releases/name=$release
  value:
    name: $release
    version: $version
    sha1: $sha1
    url: $url
EOF

bosh int git-kubo-deployment/manifests/cfcr.yml \
  -o update-$release-release.yml > git-kubo-deployment-output/manifests/cfcr.yml

pushd git-kubo-deployment-output
git config --global user.email "cfcr+cibot@pivotal.io"
git config --global user.name "CFCR CI BOT"
git add manifests/cfcr.yml
git commit -m "Bump $release to version $version"
