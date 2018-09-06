#!/usr/bin/env bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubo-version/version)
git config --global user.name "cf-london"
git config --global user.email "cf-london-eng@pivotal.io"

cp -r git-kubo-release-master/. git-kubo-release-master-output

cd git-kubo-release-master-output

cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF

bosh create-release --final --version="${version}" --sha2 --tarball "../kubo-release-tarball/kubo-release-${version}.tgz"

