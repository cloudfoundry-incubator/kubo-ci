#!/bin/bash

set -eux
set -o pipefail

cp -a kubo-lock-repo/. modified-repo
cd modified-repo
git config user.email "ci-bot@localhost"
git config user.name "CI Bot"
git mv "${POOL_NAME}/claimed/${ENV_NAME}" "${POOL_NAME}/unclaimed/${ENV_NAME}"
git commit -m "Unclaiming: ${POOL_NAME}/${ENV_NAME}"
