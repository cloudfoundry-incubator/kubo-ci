#!/bin/bash

set -euxo pipefail

source git-kubo-ci/pks-pipelines/osl/utils/all-env.sh

bosh run-errand \
  -d "${DEPLOYMENT_NAME}" \
  print-component-version \
  --json | \
  jq .Tables[0].Rows[0].stdout --raw-output \
  > osl-kubo-output.yml
echo "kubo YML output: -------"
cat osl-kubo-output.yml
echo "------------------------"

bosh run-errand \
  -d "${DEPLOYMENT_NAME}" \
  print-etcd-component-version \
  > osl-etcd-output.yml
echo "etcd YML output: -------"
cat osl-etcd-output.yml
echo "------------------------"


bosh run-errand \
  -d "${DEPLOYMENT_NAME}" \
  print-docker-component-version \
  > osl-docker-output.yml
echo "docker YML output: -----"
cat osl-docker-output.yml
echo "------------------------"
