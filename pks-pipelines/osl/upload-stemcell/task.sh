#!/bin/bash

set -euxo pipefail

# set relevant BOSH env vars
source git-kubo-ci/pks-pipelines/osl/utils/all-env.sh

pushd stemcell
  bosh upload-stemcell \
    --sha1 "$(cat sha1)" \
    "$(cat url)"?v="$(cat version)"
popd
