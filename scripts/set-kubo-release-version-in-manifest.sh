#!/usr/bin/env bash

set -eux -o pipefail

cp -r git-kubo-deployment/. git-kubo-deployment-output

release_version="$(cat kubo-version/version)"
sha="$(shasum kubo-release-tarball/kubo-release-${release_version}.tgz | cut -d ' ' -f 1)"
url="https://github.com/cloudfoundry-incubator/kubo-release/releases/download/v${release_version}/kubo-release-${release_version}.tgz"

cat << EOF > replace-kubo-version.yml
- type: replace
  path: /0/value/name=kubo
  value:
    name: kubo
    version: ((release_version))
    sha1: ((sha))
    url: ((url))
EOF

bosh int git-kubo-deployment/manifests/ops-files/non-precompiled-releases.yml \
  -o replace-kubo-version.yml \
  -v release_version="$release_version" \
  -v sha="$sha" \
  -v url="$url" \
  > git-kubo-deployment-output/manifests/ops-files/non-precompiled-releases.yml

cat << EOF > windows-replace-kubo-version.yml
- type: replace
  path: /0/value
  value:
    name: kubo
    version: ((release_version))
    sha1: ((sha))
    url: ((url))
EOF

sha_windows="$(shasum kubo-release-windows-tarball/kubo-release-windows-${release_version}.tgz | cut -d ' ' -f 1)"
url_windows="https://github.com/cloudfoundry-incubator/kubo-release-windows/releases/download/v${release_version}/kubo-release-windows-${release_version}.tgz"

bosh int git-kubo-deployment/manifests/ops-files/windows/add-worker.yml \
  -o windows-replace-kubo-version.yml \
  -v release_version="$release_version" \
  -v sha="$sha_windows" \
  -v url="$url_windows" \
  > git-kubo-deployment-output/manifests/ops-files/windows/add-worker.yml

git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"
cd git-kubo-deployment-output

git add .
git commit -m "Final release for v${release_version}"
git tag -a "v${release_version}" -m "Tag for version v${release_version}"
