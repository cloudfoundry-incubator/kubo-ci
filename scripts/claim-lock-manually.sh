#!/bin/bash

set -eu
set -o pipefail

cp -a kubo-lock-repo/. modified-repo
cd modified-repo
git config user.email "ci-bot@localhost"
git config user.name "CI Bot"

if git mv "$LOCK_DIR/unclaimed/$ENV_NAME" "$LOCK_DIR/claimed/$ENV_NAME"; then
  git commit -m "claiming: $LOCK_DIR/$ENV_NAME"
fi
