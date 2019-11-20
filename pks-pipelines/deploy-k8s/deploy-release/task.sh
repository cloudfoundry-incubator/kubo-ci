#!/bin/bash

set -eux pipefail

source git-kubo-ci/pks-pipelines/deploy-k8s/utils/all-env.sh

bosh --non-interactive deploy -d "${DEPLOYMENT_NAME}" \
  --var=deployment-name="${DEPLOYMENT_NAME}" \
  --var=kubo-version="${KUBO_GIT_SHA}" \
  --var=etcd-version="${ETCD_GIT_SHA}" \
  --var=docker-version="${DOCKER_GIT_SHA}" \
  git-kubo-ci/pks-pipelines/manifest.yml