#!/usr/bin/env bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubo-version/version)
git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"

echo "kubo-release ${version}" >kubo-release-tarball-notes/name
echo "See [CFCR Release notes](https://docs-cfcr.cfapps.io/overview/release-notes/) page" > kubo-release-tarball-notes/body

cp -r git-kubo-release/. git-kubo-release-output

cd git-kubo-release-output

cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF

release_name=${release_name:-kubo-release}

bosh create-release --final --version="${version}" --sha2 --tarball "../kubo-release-tarball/${release_name}-${version}.tgz"

git add .
git commit -m "Final release for v${version}"
git tag -a "v${version}" -m "Tag for version v${version}"
