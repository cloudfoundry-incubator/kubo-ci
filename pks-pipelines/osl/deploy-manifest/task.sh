#!/bin/bash

set -euxo pipefail

source git-kubo-ci/pks-pipelines/osl/utils/all-env.sh

bosh deploy \
  --non-interactive \
  --deployment="${DEPLOYMENT_NAME}" \
  --var=deployment-name="${DEPLOYMENT_NAME}" \
  --var=kubo-version="${KUBO_GIT_SHA}" \
  --var=kubo-windows-version="${KUBO_WINDOWS_GIT_SHA}" \
  --var=etcd-version="${ETCD_GIT_SHA}" \
  --var=docker-version="${DOCKER_GIT_SHA}" \
  git-kubo-ci/pks-pipelines/osl/manifest.yml
