#!/bin/bash

set -eux pipefail

source git-kubo-ci/pks-pipelines/deploy-k8s/utils/all-env.sh

bosh delete-deployment \
  --non-interactive \
  --deployment="${DEPLOYMENT_NAME}"
bosh delete-release \
  --non-interactive \
  kubo/"${KUBO_GIT_SHA}"
