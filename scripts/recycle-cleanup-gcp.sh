#!/usr/bin/env bash
set -eux -o pipefail

ENV_FILE=./metadata ./git-kubo-ci/scripts/cleanup-environment.sh