#!/bin/bash

set -exu -o pipefail

pool_name=$(basename $(cd "$(dirname "$ENV_FILE")/.."; pwd))
env_name=$(basename "$ENV_FILE")

cp -a kubo-lock-repo/. modified-repo

cd modified-repo
git config user.email "ci-bot@localhost"
git config user.name "CI Bot"
git mv "${pool_name}/claimed/${env_name}" "${pool_name}/unclaimed/${env_name}"
git commit -m "Unclaiming: ${pool_name}/${env_name}"
