#!/bin/bash

set -euxo pipefail

source git-kubo-ci/pks-pipelines/minimum-release-verification/utils/all-env.sh

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

# stdout will be 'Kubernetes v1.15.5'
# cut will trim that down to 1.15.5
WINDOWS_VERSION="$(cat windows_versions.txt | \
  jq .Tables[0].Rows[0].stdout --raw-output | \
  cut -d'v' -f2)"

# stdout will be OSM format yml, such as:
# other:kubernetes:1.15.5:
#  name: kubernetes
#  version: 1.15.5
#  repository: Other
#  license: Apache2.0
#  other-distribution: /tmp/osl/v1.15.5.tar.gz
#  url: https://github.com/kubernetes/kubernetes/archive/v1.15.5.tar.gz
# grep finds the header line
# cut trims down to 1.15.5
LINUX_VERSION="$(cat linux_versions.txt | \
  jq .Tables[0].Rows[0].stdout --raw-output | \
  grep '^other:kubernetes:' | \
  cut -d':' -f3)"

if [ "$WINDOWS_VERSION" == "$LINUX_VERSION" ]
then
  echo "Versions match! $LINUX_VERSION"
else
  echo "Version mismatch!  Linux: $LINUX_VERSION, Windows: $WINDOWS_VERSION"
  exit 1
fi
