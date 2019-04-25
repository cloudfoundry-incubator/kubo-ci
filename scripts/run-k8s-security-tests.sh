#!/usr/bin/env bash

set -eu -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

export GOPATH="${ROOT}/git-kubo-ci"
DEPLOYMENT_NAME="${DEPLOYMENT_NAME:="ci-service"}"

target_bosh_director() {
  if [[ -f source-json/source.json ]]; then
    source="source-json/source.json"
  else
    source="kubo-lock/metadata"
    DEPLOYMENT_NAME="$(bosh int kubo-lock/metadata --path=/deployment_name)"
  fi
  export DEPLOYMENT_NAME="${DEPLOYMENT_NAME}"
  source "${ROOT}/git-kubo-ci/scripts/set-bosh-env" ${source}
}

main() {
  if bosh int kubo-lock/metadata --path=/jumpbox_ssh_key &>/dev/null ; then
    bosh int kubo-lock/metadata --path=/jumpbox_ssh_key > ssh.key
    chmod 0600 ssh.key
    cidr="$(bosh int kubo-lock/metadata --path=/internal_cidr)"
    jumpbox_url="$(bosh int kubo-lock/metadata --path=/jumpbox_url)"
    sshuttle -r "jumpbox@${jumpbox_url}" "${cidr}" -e "ssh -i ssh.key -o StrictHostKeyChecking=no -o ServerAliveInterval=300 -o ServerAliveCountMax=10" --daemon
    trap 'kill -9 $(cat sshuttle.pid)' EXIT
  fi

  target_bosh_director

  ginkgo -keepGoing -r -progress -flakeAttempts=2 "${ROOT}/git-kubo-ci/src/tests/security-tests/"
}

main
