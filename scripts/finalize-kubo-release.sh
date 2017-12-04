#!/usr/bin/env bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubo-version/version)

pushd git-kubo-release

cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF

bosh-cli finalize-release ../gcs-kubo-release-tarball/kubo-release-*.tgz --version=${version}

git add .
git commit -m "Final release for ${version}"
git tag -a ${version} -m "Tagging for version ${version}"
