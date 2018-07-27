#!/bin/bash

set -exu -o pipefail

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"

release_tarball=$(find ${ROOT}/gcs-kubo-release-tarball/ -name "*kubo-*.tgz" | head -n1)
bosh upload-release "$release_tarball"

pushd "$ROOT"
"${ROOT}/${BOSH_DEPLOY_COMMAND}"
popd
