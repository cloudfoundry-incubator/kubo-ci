#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

if [[ $# -lt 3 ]]; then
    echo "Usage:" >&2
    echo "$0 GIT_KUBO_DEPLOYMENT_DIR DEPLOYMENT_NAME KUBO_ENVIRONMENT_DIR" >&2
    exit 1
fi

function call_bosh {
  director_ip=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip")
  BOSH_CA_CERT=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/default_ca/ca") bosh-cli -e "$director_ip" "$@"
}

GIT_KUBO_DEPLOYMENT_DIR=$1
DEPLOYMENT_NAME=$2
KUBO_ENVIRONMENT_DIR=$3

"$GIT_KUBO_DEPLOYMENT_DIR/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}"

routing_mode=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing_mode")
director_name=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name")
GIT_KUBO_CI=$(cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
GOPATH="$GIT_KUBO_CI"
export GOPATH

export PATH_TO_KUBECONFIG="$HOME/.kube/config"
TLS_KUBERNETES_CERT=$(bosh-cli int <(credhub get -n "${director_name}/${DEPLOYMENT_NAME}/tls-kubernetes" --output-json) --path='/value/certificate')
TLS_KUBERNETES_PRIVATE_KEY=$(bosh-cli int <(credhub get -n "${director_name}/${DEPLOYMENT_NAME}/tls-kubernetes" --output-json) --path='/value/private_key')
export TLS_KUBERNETES_CERT TLS_KUBERNETES_PRIVATE_KEY

if [[ ${routing_mode} == "cf" ]]; then
  KUBERNETES_SERVICE_HOST=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/cf-tcp-router-name")
  KUBERNETES_SERVICE_PORT=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/external-kubo-port")
  WORKLOAD_TCP_PORT=$(expr "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/external-kubo-port")" + 1000)
  INGRESS_CONTROLLER_TCP_PORT=$(expr "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/external-kubo-port")" + 2000)
  TCP_ROUTER_DNS_NAME=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/cf-tcp-router-name")
  CF_APPS_DOMAIN=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing-cf-app-domain-name")
  export KUBERNETES_SERVICE_HOST KUBERNETES_SERVICE_PORT WORKLOAD_TCP_PORT INGRESS_CONTROLLER_TCP_PORT TCP_ROUTER_DNS_NAME CF_APPS_DOMAIN

  ginkgo "$GOPATH/src/integration-tests/cloudfoundry"
elif [[ ${routing_mode} == "iaas" ]]; then
  WORKLOAD_ADDRESS=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_worker_host")
  WORKLOAD_PORT=$(bosh-cli int "${PWD}/git-kubo-ci/specs/nginx.yml" --path="/spec/ports/0/nodePort")
  export WORKLOAD_ADDRESS WORKLOAD_PORT

  ginkgo "$GOPATH/src/integration-tests/workload"
elif [[ ${routing_mode} == "proxy" ]]; then
  WORKLOAD_ADDRESS=$(call_bosh -d "${DEPLOYMENT_NAME}" vms --column=Instance --column=IPs | grep 'worker/' | head -1 | awk '{print $2}')
  WORKLOAD_PORT=$(bosh-cli int "${GIT_KUBO_CI}/specs/nginx.yml" --path="/spec/ports/0/nodePort")
  export WORKLOAD_ADDRESS WORKLOAD_PORT

  ginkgo "$GOPATH/src/integration-tests/workload"
fi
