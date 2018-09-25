#!/usr/bin/env bash

set -eux -o pipefail

cp -r git-kubo-deployment/. git-kubo-deployment-output

cat << EOF > replace-kubo-version.yml
- type: replace
  path: /releases/name=kubo
  value:
    name: kubo
    version: ((version))
    sha1: ((sha))
    url: ((url))
EOF
version="$(cat kubo-version/version)"
sha="$(shasum kubo-release-tarball/kubo-release-${version}.tgz | cut -d ' ' -f 1)"
url="https://github.com/cloudfoundry-incubator/kubo-release/releases/download/v${version}/kubo-release-${version}.tgz"
bosh int git-kubo-deployment/manifests/cfcr.yml -o replace-kubo-version.yml -v version="$version" -v sha="$sha" -v url="$url" > git-kubo-deployment-output/manifests/cfcr.yml

git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"
cd git-kubo-deployment-output

git add .
git commit -m "Final release for v${version}"
git tag -a "v${version}" -m "Tag for version v${version}"
