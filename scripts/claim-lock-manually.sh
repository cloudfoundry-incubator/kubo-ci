#!/bin/bash

set -exu -o pipefail
 
cp -a kubo-lock-repo/. modified-repo
cd modified-repo
git config user.email "ci-bot@localhost"
git config user.name "CI Bot"

if git mv "${POOL_NAME}/unclaimed/${ENV_NAME}" "${POOL_NAME}/claimed/${ENV_NAME}"; then
  git commit -m "Claiming: ${POOL_NAME}/${ENV_NAME}"
fi
