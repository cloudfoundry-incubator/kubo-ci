#!/bin/bash

set -exu -o pipefail

tag=$(cat "$PWD/flannel-release/tag")
version=$(cat "$PWD/flannel-release/version")
name="flannel-${tag}-linux-amd64.tar.gz"
flannel_blob_path="$PWD/flannel-release/${name}"

cp -r flannel-release/tag flannel-tag
cp -r git-kubo-release/. git-kubo-release-output

cd git-kubo-release-output

cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF

existing_flannel_spec=$(bosh blobs | grep flannel | awk '{print $1}')

bosh remove-blob $(bosh blobs | grep flannel | awk '{print $1}')
bosh add-blob ${flannel_blob_path} ${name}
bosh upload-blobs

pushd packages/flanneld

sed -E -i -e "s/([0-9]+\.)+[0-9]+/${version}/" packaging
sed -E -i -e "s/${existing_flannel_spec}/${name}/" spec

popd 

