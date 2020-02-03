#!/bin/bash

set -euxo pipefail

# define KUBO_WINDOWS_GIT_SHA as latest commit
source git-kubo-ci/pks-pipelines/osl/utils/all-env.sh

pushd git-pks-kubo-release-windows
  bosh create-release --version="${KUBO_WINDOWS_GIT_SHA}" --tarball pipeline.tgz
  bosh upload-release pipeline.tgz
popd
