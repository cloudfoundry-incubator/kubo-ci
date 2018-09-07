#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"

export GOPATH="${ROOT}/git-kubo-ci"

main() {
  pushd "${ROOT}/bbr-cli"
    tar xvf bbr-*.tar
  popd

  BOSH_ENVIRONMENT="$(jq -r '.target' "${ROOT}/source-json/source.json")"
  BOSH_CLIENT="$(jq -r '.client' "${ROOT}/source-json/source.json")"
  BOSH_CLIENT_SECRET="$(jq -r '.client_secret' "${ROOT}/source-json/source.json")"
  BOSH_CA_CERT="$(jq -r '.ca_cert' "${ROOT}/source-json/source.json")"
  BOSH_DEPLOYMENT="$DEPLOYMENT_NAME"
  KUBECONFIG="${ROOT}/gcs-kubeconfig/config"
  PATH="$PATH:${ROOT}/bbr-cli/releases/"

  export BOSH_ENVIRONMENT BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_CA_CERT BOSH_DEPLOYMENT KUBECONFIG

  ginkgo -r -progress "${ROOT}/git-kubo-ci/src/tests/bbr-tests/"
}

main
