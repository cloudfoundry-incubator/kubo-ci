#!/bin/bash
bosh int kubo-lock/metadata --path=/jumpbox_ssh_key > ssh.key
chmod 0600 ssh.key
cidr="$(bosh int kubo-lock/metadata --path=/internal_cidr)"
jumpbox_url="$(bosh int kubo-lock/metadata --path=/jumpbox_url)"
sshuttle -r "jumpbox@${jumpbox_url}" "${cidr}" -e "ssh -i ssh.key -o StrictHostKeyChecking=no -o ServerAliveInterval=300 -o ServerAliveCountMax=10" --daemon
trap 'kill -9 $(cat sshuttle.pid)' EXIT

BOSH_DEPLOYMENT="${DEPLOYMENT_NAME}"
BOSH_ENVIRONMENT=$(bosh int kubo-lock/metadata --path '/target')
BOSH_CLIENT=$(bosh int kubo-lock/metadata --path '/client')
BOSH_CLIENT_SECRET=$(bosh int kubo-lock/metadata --path '/client_secret')
BOSH_CA_CERT=$(bosh int kubo-lock/metadata --path '/ca_cert')
export BOSH_DEPLOYMENT BOSH_ENVIRONMENT BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_CA_CERT

echo "${BOSH_CA_CERT}" > bosh_ca.tmp
trap 'rm bosh_ca.tmp' EXIT

bosh -n -d turbulence deploy ./git-turbulence-release/manifests/example.yml \
  -o git-kubo-ci/manifests/ops-files/use-dev-turbulence.yml \
  -o git-kubo-ci/manifests/ops-files/use-xenial-stemcell.yml \
  -v turbulence_api_ip="10.0.255.0" \
  -v director_ip=$(bosh int "kubo-lock/metadata" --path=/internal_ip) \
  -v director_client=${BOSH_CLIENT} \
  -v director_client_secret=${BOSH_CLIENT_SECRET} \
  --var-file director_ssl.ca=bosh_ca.tmp

