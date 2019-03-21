#!/bin/bash

set -exu -o pipefail

cp -r git-kubo-deployment/. git-kubo-deployment-output

release="${RELEASE_NAME}"
version="$(cat boshrelease/version)"
opsfile="${OPSFILE:-non-precompiled-release.yml}"
mode="${MODE}"

url="$(cat boshrelease/url)"

sha1="$(cat boshrelease/sha1)"

if [ "$mode" == "name" ]; then
  cat > update-$release-release.yml <<EOF
- type: replace
  path: /name=$release/value
  value:
    name: $release
    version: $version
    sha1: $sha1
    url: $url
EOF
else
  cat > update-$release-release.yml <<EOF
- type: replace
  path: /0/value/name=$release
  value:
    name: $release
    version: $version
    sha1: $sha1
    url: $url
EOF
fi

bosh int "git-kubo-deployment/manifests/ops-files/$opsfile" \
  -o update-$release-release.yml > "git-kubo-deployment-output/manifests/ops-files/$opsfile"

pushd git-kubo-deployment-output
git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"

if [ -n "$(git status --porcelain)" ]; then
    git add "manifests/ops-files/$opsfile"
    git commit -m "Bump $release to version $version in $opsfile"
else
    echo "No changes to commit."
fi
