#!/usr/bin/env bash

set -exu -o pipefail

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
version=$(cat kubo-version/version)

pushd git-kubo-release
cat .git/ref #TODO remove debugging
cp .git/ref ../kubo-release-git-ref/kubo-release-git-ref