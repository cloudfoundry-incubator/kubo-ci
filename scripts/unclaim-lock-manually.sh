#!/bin/bash

set -eu
set -o pipefail

cp -a kubo-lock-repo/. modified-repo
cd modified-repo
git config user.email "ci-bot@localhost"
git config user.name "CI Bot"
git mv "$LOCK_DIR/claimed/$ENV_NAME" "$LOCK_DIR/unclaimed/$ENV_NAME"
git commit -m "unclaiming: $LOCK_DIR/$ENV_NAME"
