#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

TURBULENCE_USERNAME="turbulence"
TURBULENCE_PASSWORD=$(bosh-cli int "$PWD/gcs-bosh-creds/creds.yml" --path='/turbulence_api_password')
bosh_ip=$(bosh-cli int "$PWD/kubo-lock/metadata" --path='/internal_ip')
TURBULENCE_API_ENDPOINT="$bosh_ip:8080/api/v1"

GIT_KUBO_CI=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)
GOPATH="$GIT_KUBO_CI"
export GOPATH


export TURBULENCE_USERNAME
export TURBULENCE_PASSWORD
export TURBULENCE_API_ENDPOINT
ginkgo "$GOPATH/src/turbulence-tests/worker_failure"

# As of now we are not testing for the worker to come back up
# We have to wait for the worker to be resurrected
# before tearing down the deployment
sleep 600