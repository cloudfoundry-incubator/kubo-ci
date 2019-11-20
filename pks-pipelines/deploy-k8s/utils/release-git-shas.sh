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

#TODO: get these from git repos
DOCKER_GIT_SHA="35.3.4"

export KUBO_GIT_SHA ETCD_GIT_SHA DOCKER_GIT_SHA
