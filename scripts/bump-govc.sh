#!/bin/bash

set -exu -o pipefail

tag=$(cat "$PWD/govc-release/tag")
version=$(cat "$PWD/govc-release/version")
name="govc_${tag}_linux_amd64.gz"
govc_blob_path="$PWD/govc-release/govc_linux_amd64.gz"

cp -r git-kubo-release/. git-kubo-release-output

cd git-kubo-release-output

cat <<EOF > "config/private.yml"
blobstore:
  options:
    json_key: '${GCS_JSON_KEY}'
EOF

existing_govc_spec=$(bosh blobs | grep govc | awk '{print $1}')

if [ $name == $existing_govc_spec ]; then
  echo "Govc blob already up-to-date."
  exit 0
fi

bosh remove-blob ${existing_govc_spec}
bosh add-blob ${govc_blob_path} ${name}
bosh upload-blobs

pushd packages/govc
sed -E -i -e "s/([0-9]+\.)+[0-9]+/${version}/" packaging
sed -E -i -e "s/${existing_govc_spec}/${name}/" spec
popd

git config --global user.email "cfcr+cibot@pivotal.io"
git config --global user.name "CFCR CI BOT"
git add .
git commit -m "Bump govc to $tag"
