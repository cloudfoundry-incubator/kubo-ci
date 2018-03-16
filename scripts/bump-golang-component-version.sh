#!/bin/bash

set -exu -o pipefail

existing_component_version=$(cat PWD/golang-version/version)

golang_rel_version=$(cat PWD/git-golang-release/version)
golang_rel_tag=$(cat PWD/git-golang-release/tag)

# Extract source tarball to a directory
mkdir golang-release-tarball
tar -zxvf "git-golang-release/${golang_rel_tag}.tar.gz" -C golang-release-tarball/

pushd golang-release-tarball
versions=$(spruce json config/blobs.yml | jq -r "keys[] | select(match(\".linux-amd64.tar.gz\"))" --raw-output | grep -oE "([0-9]+\.)+[0-9]+" | awk '{print $1}')
golang_semvers=(${versions// /\n})
popd

# Since we know that golang-release has two go versions, we find the greatest of those two
latest_golang_version=${golang_semvers[0]}
if [ $(semver ${golang_semvers[0]} ${golang_semvers[1]}) -lt 0 ]; then
    latest_golang_version=${arr[1]}
fi

# Now we compare the existing version with the latest version
if [ $(semver $latest_golang_version $existing_component_version) -le 0 ]; then
    echo "existing golang component version is already at the latest version..."
    exit 0
fi

cp -a golang-version/. modified-golang-version
