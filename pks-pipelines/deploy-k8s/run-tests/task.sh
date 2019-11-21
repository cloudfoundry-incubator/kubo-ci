#!/bin/bash

set -eux pipefail

source git-kubo-ci/pks-pipelines/deploy-k8s/utils/all-env.sh

bosh run-errand \
  -d "${DEPLOYMENT_NAME}" \
  print-component-version \
  --json \
  > linux_versions.txt

bosh run-errand \
  -d "${DEPLOYMENT_NAME}" \
  print-kubo-windows-component-version \
  --json \
  > windows_versions.txt

exit 1
