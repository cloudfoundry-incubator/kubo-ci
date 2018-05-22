#!/bin/bash

set -exu -o pipefail

source git-kubo-ci/scripts/lib/semver.sh

HOME_DIR="$PWD"
GO_VERSION=$(cat $PWD/golang-version/component-golang-version)
EXISTING_V=""

check_and_remove_existing_vendor_package() {
  pushd "$HOME_DIR"/modified-release/ > /dev/null

  EXISTING_V=$(ls -al packages | grep golang | grep -oE "([0-9]+\.)+[0-9]+")
  if [ $(compare_semvers $GO_VERSION $EXISTING_V) -eq 0 ]; then
    echo "Release already at the latest golang vendor package"
    exit 0
  fi
  pushd .final_builds/packages >/dev/null; rm -rf "golang-${EXISTING_V}-linux/"; popd >/dev/null;
  pushd packages >/dev/null; rm -rf "golang-${EXISTING_V}-linux/"; popd >/dev/null;

  popd > /dev/null
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
    local GO_VERSION_PARTS PACKAGE_GOLANG_VERSION
    semver_arr "${GO_VERSION}" GO_VERSION_PARTS
    PACKAGE_GO_VERSION="${GO_VERSION_PARTS[0]}.${GO_VERSION_PARTS[1]}"
    bosh vendor-package golang-"${PACKAGE_GO_VERSION}"-linux "$HOME_DIR"/golang-release

    grep --exclude=spec.lock --exclude-dir=src --exclude-dir=.git --exclude-dir=releases -r -o -l -E 'golang-([0-9]+\.)+[0-9]+' | xargs sed -E -i -e "/golang/s/([0-9]+\.)+[0-9]+/${PACKAGE_GO_VERSION}/"

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
