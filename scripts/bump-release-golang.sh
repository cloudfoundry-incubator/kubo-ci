#!/bin/bash

set -exu -o pipefail

MINOR_GO_VERSION="1.9"
HOME_DIR="$PWD"

extract_golang_release() {
  mkdir golang
  tar -xzvf golang-release/*.tar.gz -C "$HOME_DIR"/golang
}

vendor_golang() {
  pushd modified-release
    bosh vendor-package golang-"$MINOR_GO_VERSION"-linux "$HOME_DIR"/golang/bosh-packages-golang-release-*
    go_version="$(bosh int packages/golang-"$MINOR_GO_VERSION"-linux/spec.lock --path=/name)"

    git config --global user.email "cfcr+cibot@pivotal.io"
    git config --global user.name "CFCR CI BOT"
    git add .
    git commit -m "Updates golang to version $go_version"
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
