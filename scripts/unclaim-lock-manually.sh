#!/bin/bash

set -eu
set -o pipefail

cp -a kubo-lock-repo/. modified-repo
cd modified-repo
git config user.email "ci-bot@localhost"
git config user.name "CI Bot"
git mv "kubo-deployment/claimed/$ENV_NAME" "kubo-deployment/unclaimed/$ENV_NAME"
git commit -m "unclaiming: $ENV_NAME"
