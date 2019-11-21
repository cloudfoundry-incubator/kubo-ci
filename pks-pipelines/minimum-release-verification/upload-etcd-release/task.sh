#!/bin/bash

set -eux pipefail

source git-kubo-ci/pks-pipelines/minimum-release-verification/utils/all-env.sh

pushd git-pks-cfcr-etcd-release
  bosh create-release --version="${ETCD_GIT_SHA}" --tarball pipeline.tgz
  bosh upload-release pipeline.tgz
popd
