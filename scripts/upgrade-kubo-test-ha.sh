#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"

update_kubo() {
  echo "Updating Kubo..."
  ${DIR}/deploy-k8s-instance.sh
}

BOSH_ENV="$KUBO_ENVIRONMENT_DIR" source "$KUBO_DEPLOYMENT_DIR/bin/set_bosh_environment"
bosh-cli upload-release https://bosh.io/d/github.com/cf-platform-eng/docker-boshrelease?v=28.0.1 --sha1 448eaa2f478dc8794933781b478fae02aa44ed6b
bosh-cli upload-release https://github.com/pivotal-cf-experimental/kubo-etcd/releases/download/v2/kubo-etcd.2.tgz --sha1 ae95e661cd9df3bdc59ee38bf94dd98e2f280d4f

run_upgrade_test update_kubo
