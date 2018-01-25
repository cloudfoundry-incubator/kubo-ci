#!/bin/bash

set -exu -o pipefail

MINOR_GO_VERSION="1.9"
HOME_DIR="$PWD"

extract_golang_release() {
  mkdir golang
  tar -xzvf golang-release/*.tar.gz -C "$HOME_DIR"/golang
}

vendor_golang() {
  pushd "$HOME_DIR"/golang/bosh-packages-golang-release-*
    blob_name="$(bosh blobs | grep linux | grep 1.9 | awk '{print $1}')"
    go_version="${blob_name%.tar.gz}"
  popd

  pushd modified-release
    cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF
    bosh vendor-package golang-"$MINOR_GO_VERSION"-linux "$HOME_DIR"/golang/bosh-packages-golang-release-*

    git config --global user.email "cfcr+cibot@pivotal.io"
    git config --global user.name "CFCR CI BOT"

    set +e
    git add .
    git commit -m "Updates golang to version $go_version"
    set -e
  popd
}

create_output_directory() {
  cp -a release/. modified-release
}

main() {
  create_output_directory
  extract_golang_release
  vendor_golang
}

main
