#!/bin/bash

set -eux pipefail

source git-kubo-ci/pks-pipelines/deploy-k8s/utils/lock-to-env.sh
source git-kubo-ci/pks-pipelines/deploy-k8s/utils/release-git-shas.sh
