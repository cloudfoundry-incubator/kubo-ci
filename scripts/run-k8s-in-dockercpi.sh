#!/usr/bin/env bash

# fly -t kubo execute -p --config tasks/bosh-docker-cpi.yml \
#  -i git-kubo-ci=./ -i git-kubo-deployment=../kubo-deployment

set -eu

start-bosh -o /usr/local/bosh-deployment/uaa.yml -o /usr/local/bosh-deployment/credhub.yml
source /tmp/local-bosh/director/env

bosh cloud-config || true

bosh runtime-config || true

bosh releases
bosh upload-stemcell https://bosh.io/d/stemcells/bosh-warden-boshlite-ubuntu-trusty-go_agent
bosh upload-stemcell https://bosh.io/d/stemcells/bosh-warden-boshlite-ubuntu-xenial-go_agent
bosh upload-release https://bosh.io/d/github.com/cloudfoundry-incubator/kubo-release

bosh -n -d cfcr deploy --no-redact \
  git-kubo-deployment/manifests/cfcr.yml \
  -o git-kubo-deployment/manifests/ops-files/misc/single-master.yml \
  "$@"

bosh -d cfcr run-errand apply-specs

bosh -d cfcr run-errand smoke-tests
