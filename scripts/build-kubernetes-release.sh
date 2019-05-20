#!/bin/bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubernetes-version/version)

cd git-kubernetes-release

bosh create-release --name "kubernetes" --sha2 --tarball="../kubernetes-release/kubernetes-release-${version}.tgz" --version=${version}
