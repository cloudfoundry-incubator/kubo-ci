#!/bin/bash

set -euxo pipefail

source git-kubo-ci/pks-pipelines/osl/utils/all-env.sh

#####################
# Delete Deployment #
#####################
bosh delete-deployment \
  --non-interactive \
  --deployment="${DEPLOYMENT_NAME}"

###################
# Delete Releases #
###################
bosh delete-release \
  --non-interactive \
  kubo/"${KUBO_GIT_SHA}" \
  || true
bosh delete-release \
  --non-interactive \
  cfcr-etcd/"${ETCD_GIT_SHA}" \
  || true
bosh delete-release \
  --non-interactive \
  docker/"${DOCKER_GIT_SHA}" \
  || true
