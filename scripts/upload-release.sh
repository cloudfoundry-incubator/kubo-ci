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
  bosh upload-release https://storage.googleapis.com/kubo-public/docker-35.2.3-ubuntu-xenial-315.36-20190716-163114-008878.tgz
  bosh upload-release https://storage.googleapis.com/kubo-public/docker-35.2.3-windows2019-2019.7-20190716-161813-432556.tgz
  bosh upload-release https://storage.googleapis.com/kubo-public/kubo-1.0.0-dev.5-ubuntu-xenial-315.70-20190719-022102-888261013.tgz
  bosh upload-release https://storage.googleapis.com/kubo-public/kubo-1.0.0-dev.5-windows2019-2019.4-20190719-022425-676949049.tgz
}

