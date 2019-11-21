#!/bin/bash

set -eux pipefail

source git-kubo-ci/scripts/set-bosh-env kubo-lock/metadata

DEPLOYMENT_NAME="$(bosh interpolate kubo-lock/metadata --path=/deployment_name)"
export DEPLOYMENT_NAME
