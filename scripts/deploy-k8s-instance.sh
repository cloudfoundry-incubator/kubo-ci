#!/bin/bash

set -exu -o pipefail

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )"

pushd "$ROOT"
"${ROOT}/${BOSH_DEPLOY_COMMAND}"
popd
