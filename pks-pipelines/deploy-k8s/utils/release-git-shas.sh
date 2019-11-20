#!/bin/bash

set -eux pipefail

pushd git-pks-kubo-release
  KUBO_GIT_SHA="$(git log -1 --format='%H')"
popd

pushd git-pks-cfcr-etcd-release
  ETCD_GIT_SHA="$(git log -1 --format='%H')"
popd

#TODO: get these from git repos
DOCKER_GIT_SHA="35.3.4"

export KUBO_GIT_SHA ETCD_GIT_SHA DOCKER_GIT_SHA
