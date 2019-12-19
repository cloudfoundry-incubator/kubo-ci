#!/bin/bash

set -euxo pipefail

# set relevant BOSH env vars
source git-kubo-ci/pks-pipelines/minimum-release-verification/utils/all-env.sh

files=( stemcell/bosh-stemcell-*.tgz )
bosh upload-stemcell "${files[0]}"
