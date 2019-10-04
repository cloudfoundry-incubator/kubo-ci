#!/bin/bash

set -exu -o pipefail

release=${release:-kubo}
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubo-version/version)

cd git-kubo-release

bosh create-release --name ${release} --sha2 --tarball="../${release}-release/${release}-release-${version}.tgz" --version=${version}
