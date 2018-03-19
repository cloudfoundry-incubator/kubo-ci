#!/bin/bash

set -exu -o pipefail

source git-kubo-ci/scripts/lib/semver.sh

existing_component_version=$(cat $PWD/golang-version/component-golang-version)

# Extract source tarball to a directory
mkdir golang-release-tarball
tar -zxf "git-golang-release/source.tar.gz" -C golang-release-tarball/
mv golang-release-tarball/bosh-*/* golang-release-tarball/

pushd golang-release-tarball
versions=$(spruce json config/blobs.yml | jq -r "keys[] | select(match(\".linux-amd64.tar.gz\"))" --raw-output | grep -oE "([0-9]+\.)+[0-9]+" | awk '{print $1}')
golang_semvers=(${versions// /\n})
popd

# Since we know that golang-release has two go versions, we find the greatest of those two
latest_golang_version=${golang_semvers[0]}
if [ $(compare_semvers ${golang_semvers[0]} ${golang_semvers[1]}) -lt 0 ]; then
    latest_golang_version=${arr[1]}
fi

# Now we compare the existing version with the latest version
if [ $(compare_semvers $latest_golang_version $existing_component_version) -le 0 ]; then
    echo "existing golang component version is already at the latest version..."
    exit 0
else
    echo $latest_golang_version > golang-version/component-golang-version
    # removing the newline character at the end
    truncate -s -1 golang-version/component-golang-version
fi

cp -a golang-version/. modified-golang-version
