#!/bin/bash

set -exu -o pipefail

source git-kubo-ci/scripts/lib/semver.sh

HOME_DIR="$PWD"
GO_VERSION=$(cat $PWD/golang-version/component-golang-version)
EXISTING_V=""

check_and_remove_existing_vendor_package() {
  pushd "$HOME_DIR"/modified-release/.final_builds/packages
    EXISTING_V=$(ls -al | grep golang | grep -oE "([0-9]+\.)+[0-9]+")
    if [ $(compare_semvers $GO_VERSION $EXISTING_V) -le 0 ]; then
      echo "Release already at the latest golang vendor package"
      exit 0
    fi
    rm -rf "golang-${EXISTING_V}-linux/"
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
    git add -A
    git commit -m "Updates golang to version $go_version"
    set -e
  popd
}

create_output_directory() {
  cp -a release/. modified-release
}

output_existing_version() {
  echo $EXISTING_V > existing_golang_version
  truncate -s -1 existing_golang_version
}

main() {
  create_output_directory
  check_and_remove_existing_vendor_package
  output_existing_version
  vendor_golang
}

main
