#!/bin/bash

set -eux pipefail

source git-kubo-ci/pks-pipelines/deploy-k8s/utils/all-env.sh

pushd git-pks-kubo-release
  bosh create-release --version="${KUBO_GIT_SHA}"
  bosh upload-release
popd
