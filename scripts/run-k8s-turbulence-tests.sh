#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

TURBULENCE_USERNAME="turbulence"
TURBULENCE_PASSWORD=$(bosh-cli int "$PWD/gcs-bosh-creds/creds.yml" --path='/turbulence_api_password')
bosh_ip=$(bosh-cli int "$PWD/kubo-lock/metadata" --path='/internal_ip')
export TURBULENCE_API_ENDPOINT="$bosh_ip:8080/api/v1"

GIT_KUBO_CI=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)
GOPATH="$GIT_KUBO_CI"
export GOPATH

ginkgo "$GOPATH/src/turbulence-tests/worker_failure"
