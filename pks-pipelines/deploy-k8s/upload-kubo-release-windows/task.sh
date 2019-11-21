#!/bin/bash

set -eux pipefail

source git-kubo-ci/pks-pipelines/deploy-k8s/utils/all-env.sh

pushd git-pks-kubo-release-windows
  bosh create-release --version="${KUBO_WINDOWS_GIT_SHA}" --tarball pipeline.tgz
  bosh upload-release pipeline.tgz
popd
