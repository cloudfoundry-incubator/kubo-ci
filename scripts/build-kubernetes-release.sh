#!/bin/bash

set -exu -o pipefail

RELEASE=${RELEASE:-kubernetes}
export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubernetes-version/version)

cd git-${RELEASE}-release

bosh create-release --name "${RELEASE}" --sha2 --tarball="../${RELEASE}-release/${RELEASE}-release-${version}.tgz" --version=${version}
