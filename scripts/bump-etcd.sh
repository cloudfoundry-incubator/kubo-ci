#!/bin/bash

set -exu -o pipefail

tag=$(cat "$PWD/etcd-release/tag")
version=$(cat "$PWD/etcd-release/version")
name="etcd-${tag}-linux-amd64.tar.gz"
etcd_blob_path="$PWD/etcd-release/${name}"

cp -r git-kubo-etcd/. git-kubo-etcd-output

pushd git-kubo-etcd-output

cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF

existing_etcd_spec=$(bosh blobs | grep etcd | awk '{print $1}')

if [ $name == $existing_etcd_spec ]; then
  echo "etcd blob already up-to-date."
  exit 0
fi

bosh remove-blob ${existing_etcd_spec}
bosh add-blob ${etcd_blob_path} ${name}
bosh upload-blobs

git config --global user.email "cfcr+cibot@pivotal.io"
git config --global user.name "CFCR CI BOT"
git add .
git commit -m "Bump etcd to $tag"

popd
