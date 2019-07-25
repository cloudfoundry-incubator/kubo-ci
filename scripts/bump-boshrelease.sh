#!/bin/bash

set -exu -o pipefail

cp -r git-kubo-deployment/. git-kubo-deployment-output

release="${RELEASE_NAME}"
version="$(cat boshrelease/version)"
array_pos="${ARRAY_POS}"
base_ops_file="${BASE_OPS_FILE}"

url="$(cat boshrelease/url)"

sha1="$(cat boshrelease/sha1)"

cat > update-$release-release.yml <<EOF
- type: replace
  path: /$array_pos/value/name=$release
  value:
    name: $release
    version: $version
    sha1: $sha1
    url: $url
EOF

bosh int git-kubo-deployment/$base_ops_file \
  -o update-$release-release.yml > git-kubo-deployment-output/$base_ops_file

pushd git-kubo-deployment-output
git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"

if [ -n "$(git status --porcelain)" ]; then
    git add $base_ops_file
    git commit -m "Bump $release to version $version in $(basename $base_ops_file .yml) ops-file"
else
    echo "No changes to commit."
fi
