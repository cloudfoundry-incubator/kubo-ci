#!/bin/bash

set -eux pipefail

if [ -d "git-pks-kubo-release" ]
then
  pushd git-pks-kubo-release
    KUBO_GIT_SHA="$(git log -1 --format='%H')"
  popd
else
  KUBO_GIT_SHA=""
fi

if [ -d "git-pks-cfcr-etcd-release" ]
then
  pushd git-pks-cfcr-etcd-release
    ETCD_GIT_SHA="$(git log -1 --format='%H')"
  popd
else
  ETCD_GIT_SHA=""
fi

if [ -d "git-pks-docker-bosh-release" ]
then
  pushd git-pks-docker-bosh-release
    DOCKER_GIT_SHA="$(git log -1 --format='%H')"
  popd
else
  DOCKER_GIT_SHA=""
fi

export KUBO_GIT_SHA ETCD_GIT_SHA DOCKER_GIT_SHA
