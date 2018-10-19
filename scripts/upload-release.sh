#!/bin/bash

set -euo pipefail

target_bosh_director() {
  BOSH_ENVIRONMENT=$(bosh int source-json/source.json --path '/target')
  BOSH_CLIENT=$(bosh int source-json/source.json --path '/client')
  BOSH_CLIENT_SECRET=$(bosh int source-json/source.json --path '/client_secret')
  BOSH_CA_CERT=$(bosh int source-json/source.json --path '/ca_cert')
  export BOSH_ENVIRONMENT BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_CA_CERT
}

main() {
  target_bosh_director
  bosh upload-release ${PWD}/release/*.tgz
}

