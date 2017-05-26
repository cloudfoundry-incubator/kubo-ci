#!/bin/bash

set -eu
set -o pipefail

cp -a kubo-lock-repo/. modified-repo
cd modified-repo
git config user.email "ci-bot@localhost"
git config user.name "CI Bot"

if git mv "kubo-deployment/unclaimed/$ENV_NAME" "kubo-deployment/claimed/$ENV_NAME"; then
  git commit -m "claiming: $ENV_NAME"
fi
