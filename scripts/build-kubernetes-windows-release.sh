#!/bin/bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubernetes-windows-version/version)

cd git-kubernetes-windows-release

bosh create-release --name "kubernetes-windows" --sha2 --tarball="../kubernetes-windows-release/kubernetes-windows-release-${version}.tgz" --version=${version}
