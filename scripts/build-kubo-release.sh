#!/bin/bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubo-version/version)

cd git-kubo-release

bosh-cli create-release --name "kubo" --sha2 --tarball="../kubo-release/kubo-release-${version}.tgz" --version=${version}
