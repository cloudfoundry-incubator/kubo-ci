#!/bin/bash

set -exu -o pipefail

existing_component_version=$(cat $PWD/golang-version/component-golang-version)

pushd golang-release
latest_golang_version=$(bosh blobs --column=path | grep linux | sed "s/go\(.*\)\.linux\-amd64\.tar\.gz.*/\1/" | sort | tail -1)
popd

# Now we compare the existing version with the latest version
if [ "$latest_golang_version" == "$existing_component_version" ]; then
    echo "existing golang component version is already at the latest version..."
    exit 0
else
    echo $latest_golang_version > golang-version/component-golang-version
    # removing the newline character at the end
    truncate -s -1 golang-version/component-golang-version
fi

cp -a golang-version/. modified-golang-version
