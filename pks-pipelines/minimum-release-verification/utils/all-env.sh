#!/bin/bash

set -euxo pipefail

source git-kubo-ci/pks-pipelines/minimum-release-verification/utils/lock-to-env.sh
source git-kubo-ci/pks-pipelines/minimum-release-verification/utils/define-git-shas-of-release.sh
