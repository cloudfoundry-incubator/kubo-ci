#!/usr/bin/env bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubo-etcd-version/version | awk -F. '{print $1}')
git config --global user.name "cf-london"
git config --global user.email "cf-london-eng@pivotal.io"

cp -r git-kubo-etcd-release/. git-kubo-etcd-output

cd git-kubo-etcd-output

cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF

bosh create-release --final --version=${version} --sha2 \
  --tarball ../kubo-etcd/kubo-etcd.${version}.tgz


echo "kubo-etcd v${version}" >../kubo-etcd/name
echo "" > ../kubo-etcd/body

git checkout -b tmp/release
git add .
git commit -m "Final release for v${version}"
git tag -a "v${version}" -m "Tag for version v${version}"
git checkout master
git merge tmp/release -m "Merge release branch for v${version}"
git branch -d tmp/release
