#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

print_usage() {
  echo "Usage: VM_TYPE=[master | worker] NODE_TYPE=[master | node] run-cis-benchmarks.sh"
}

if [ -z "$VM_TYPE" ] || [ -z "$NODE_TYPE" ]; then
  print_usage
  exit 1
fi

set -eu
set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# shellcheck source=lib/environment.sh
. "$DIR/lib/environment.sh"

copy_state_and_creds() {
  DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"
  source "${DIR}/lib/utils.sh"
  setup_env "${KUBO_ENVIRONMENT_DIR}"
}

if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  copy_state_and_creds
fi

BOSH_ENV="${KUBO_ENVIRONMENT_DIR}"

BOSH_CLIENT=bosh_admin
BOSH_CLIENT_SECRET=$(bosh int "$BOSH_ENV"/creds.yml --path=/bosh_admin_client_secret)
BOSH_CA_CERT=$(bosh int "$BOSH_ENV"/creds.yml --path=/default_ca/ca)
BOSH_ENVIRONMENT=$(bosh int "$BOSH_ENV"/director.yml --path=/internal_ip)

export BOSH_NAME BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_CA_CERT BOSH_ENVIRONMENT

dst="/tmp/$(date +%s)"

bosh -d ci-service ssh "$VM_TYPE/0" --command="mkdir -p $dst"

bosh -d ci-service scp \
  "$DIR/compile-run-kube-bench.sh" "$VM_TYPE/0:$dst/compile-run-kube-bench.sh"
bosh -d ci-service scp \
  "$DIR/kube-bench/config.yml" "$VM_TYPE/0:$dst/kube-bench-config.yml"

bosh -d ci-service ssh "$VM_TYPE/0" \
  --command="cp $dst/* .; ./compile-run-kube-bench.sh $NODE_TYPE"

