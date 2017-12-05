#!/usr/bin/env bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubo-version/version)
git config --global user.name "cf-london"
git config --global user.email "cf-london-eng@pivotal.io"

cp -r git-kubo-release/. git-kubo-release-output

cd git-kubo-release-output

cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF

bosh-cli create-release --final --version=${version} --sha2 --tarball ../kubo-release/kubo-release-${version}.tgz

echo "kubo-release ${version}" >../kubo-release/name
echo "v${version}" > ../kubo-release/tag
echo "See [CFCR Release notes](https://docs-kubo.cfapps.io/overview/release-notes/) page" > ../kubo-release/body

echo "kubo-deployment ${version}" >../kubo-deployment/name
echo "v${version}" > ../kubo-deployment/tag
echo "See [CFCR Release notes](https://docs-kubo.cfapps.io/overview/release-notes/) page" > ../kubo-deployment/body

git checkout -b tmp/release
git add .
git commit -m "Final release for v${version}"
git tag -a "v${version}" -m "Tag for version v${version}"
git checkout master
git merge tmp/release -m "Merge release branch for v${version}"
git branch -d tmp/release