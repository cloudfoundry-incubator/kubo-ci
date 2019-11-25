#!/bin/bash

set -euxo pipefail

source git-kubo-ci/pks-pipelines/minimum-release-verification/utils/all-env.sh

pushd git-pks-kubo-release-windows
  bosh create-release --version="${KUBO_WINDOWS_GIT_SHA}" --tarball pipeline.tgz
  bosh upload-release pipeline.tgz
popd
