#!/bin/bash

set -exu -o pipefail

source git-kubo-ci/scripts/lib/semver.sh

HOME_DIR="$PWD"
GO_VERSION=$(cat $PWD/golang-version/component-golang-version)

check_and_remove_existing_vendor_package() {
  pushd "$HOME_DIR"/golang-release/.final_builds/packages
    local existing_v=$(ls -al | grep golang | grep -oE "([0-9]+\.)+[0-9]+")
    if [ $(compare_semvers $GO_VERSION $existing_v) -le 0 ]; then
      echo "Release ${release} already at the latest golang vendor package"
      exit 0
    fi
    rm -rf "golang-${existing_v}-linux/"
  popd
}

vendor_golang() {
  pushd "$HOME_DIR"/golang-release
    blob_name=$(bosh blobs --json | jq '.Tables[0].Rows[] | .path | select(test("'"${GO_VERSION}"'.*linux"))' --raw-output)
    go_version="${blob_name%.tar.gz}"
  popd

  pushd modified-release
    set +x
    cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF
    set -x
    bosh vendor-package golang-"${GO_VERSION}"-linux "$HOME_DIR"/golang-release

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
  check_and_remove_existing_vendor_package
  vendor_golang
}

main
