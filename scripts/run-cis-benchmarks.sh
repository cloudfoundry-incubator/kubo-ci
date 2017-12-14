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
  cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
  "$PWD/git-kubo-deployment/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "ci-service"
}

if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  copy_state_and_creds
fi

BOSH_ENV="${KUBO_ENVIRONMENT_DIR}"

BOSH_CLIENT=bosh_admin
BOSH_CLIENT_SECRET=$(bosh-cli int "$BOSH_ENV"/creds.yml --path=/bosh_admin_client_secret)
BOSH_CA_CERT=$(bosh-cli int "$BOSH_ENV"/creds.yml --path=/default_ca/ca)
BOSH_ENVIRONMENT=$(bosh-cli int "$BOSH_ENV"/director.yml --path=/internal_ip)

export BOSH_NAME BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_CA_CERT BOSH_ENVIRONMENT

dst="/tmp/$(date +%s)"

bosh-cli -d ci-service ssh "$VM_TYPE/0" --command="mkdir -p $dst"

bosh-cli -d ci-service scp \
  "$DIR/compile-run-kube-bench.sh" "$VM_TYPE/0:$dst/compile-run-kube-bench.sh"
bosh-cli -d ci-service scp \
  "$DIR/kube-bench/config.yml" "$VM_TYPE/0:$dst/kube-bench-config.yml"

bosh-cli -d ci-service ssh "$VM_TYPE/0" \
  --command="cp $dst/* .; ./compile-run-kube-bench.sh $NODE_TYPE"

