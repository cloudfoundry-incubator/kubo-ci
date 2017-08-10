#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

TURBULENCE_USERNAME="turbulence"
TURBULENCE_PASSWORD=$(bosh-cli int "$PWD/gcs-bosh-creds/creds.yml" --path='/turbulence_api_password')
BOSH_ENVIRONMENT=$(bosh-cli int "$PWD/kubo-lock/metadata" --path='/internal_ip')
BOSH_CA_CERT=$(bosh-cli int "$PWD/gcs-bosh-creds/creds.yml" --path='/default_ca/ca')
BOSH_CLIENT=bosh_admin
BOSH_CLIENT_SECRET=$(bosh-cli int "$PWD/gcs-bosh-creds/creds.yml" --path='/bosh_admin_client_secret')
TURBULENCE_API_ENDPOINT="$BOSH_ENVIRONMENT:8080/api/v1"

GIT_KUBO_CI=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)
GOPATH="$GIT_KUBO_CI"
export GOPATH


export TURBULENCE_USERNAME
export TURBULENCE_PASSWORD
export TURBULENCE_API_ENDPOINT
export BOSH_ENVIRONMENT
export BOSH_CA_CERT
export BOSH_CLIENT
export BOSH_CLIENT_SECRET
ginkgo "$GOPATH/src/turbulence-tests/worker_failure"

# As of now we are not testing for the worker to come back up
# We have to wait for the worker to be resurrected
# before tearing down the deployment
sleep 600