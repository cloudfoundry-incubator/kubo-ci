#!/bin/bash
set -exu -o pipefail

source git-kubo-ci/scripts/lib/generate-pr.sh

tag=$(cat "$PWD/k8s-release/tag")
version=$(cat "$PWD/k8s-release/version")

cp -r git-kubo-release/. git-kubo-release-output
pushd git-kubo-release-output

existing_k8s_version=$(bosh blobs --column path | grep kubelet | grep -o -E '([0-9]+\.)+[0-9]+')

if [ $existing_k8s_version == $version ]; then
    echo "Kubernetes version already up-to-date."
    exit 0
else
    ./scripts/download_k8s_binaries $version
    cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF
    bosh upload-blobs
    generate_pull_request "kubernetes" $tag
    exit 1
fi

popd
