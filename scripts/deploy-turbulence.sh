#!/bin/bash

BOSH_DEPLOYMENT="${DEPLOYMENT_NAME}"
BOSH_ENVIRONMENT=$(bosh int kubo-lock/metadata --path '/target')
BOSH_CLIENT=$(bosh int kubo-lock/metadata --path '/client')
BOSH_CLIENT_SECRET=$(bosh int kubo-lock/metadata --path '/client_secret')
BOSH_CA_CERT=$(bosh int kubo-lock/metadata --path '/ca_cert')
export BOSH_DEPLOYMENT BOSH_ENVIRONMENT BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_CA_CERT

bosh -n -d turbulence deploy ./git-turbulence-release/manifests/example.yml \
  -v turbulence_api_ip="10.100.0.10" \
  -v director_ip=$(bosh int "${ROOT}/kubo-lock/metadata" --path=/internal_ip) \
  -v director_client=${BOSH_CLIENT} \
  -v director_client_secret=${BOSH_CLIENT_SECRET} \
  --vars-store /tmp/turbulence.yml
