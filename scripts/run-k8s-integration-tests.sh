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
  BOSH_ENV="$KUBO_ENVIRONMENT_DIR" source "$GIT_KUBO_DEPLOYMENT_DIR/bin/set_bosh_environment"
  bosh-cli "$@"
}

GIT_KUBO_DEPLOYMENT_DIR=$1
DEPLOYMENT_NAME=$2
KUBO_ENVIRONMENT_DIR=$3

"$GIT_KUBO_DEPLOYMENT_DIR/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}"

routing_mode=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing_mode")
iaas=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/iaas")
director_name=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name")
GIT_KUBO_CI=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)
GOPATH="$GIT_KUBO_CI"
export GOPATH

export PATH_TO_KUBECONFIG="$HOME/.kube/config"
TLS_KUBERNETES_CERT=$(bosh-cli int <(credhub get -n "${director_name}/${DEPLOYMENT_NAME}/tls-kubernetes" --output-json) --path='/value/certificate')
TLS_KUBERNETES_PRIVATE_KEY=$(bosh-cli int <(credhub get -n "${director_name}/${DEPLOYMENT_NAME}/tls-kubernetes" --output-json) --path='/value/private_key')
export TLS_KUBERNETES_CERT TLS_KUBERNETES_PRIVATE_KEY

if [[ ${routing_mode} == "cf" ]]; then
  KUBERNETES_SERVICE_HOST=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_host")
  KUBERNETES_SERVICE_PORT=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_port")
  WORKLOAD_TCP_PORT=$(expr "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_port")" + 1000)
  INGRESS_CONTROLLER_TCP_PORT=$(expr "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_port")" + 2000)
  TCP_ROUTER_DNS_NAME=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_host")
  CF_APPS_DOMAIN=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing-cf-app-domain-name")
  export KUBERNETES_SERVICE_HOST KUBERNETES_SERVICE_PORT WORKLOAD_TCP_PORT INGRESS_CONTROLLER_TCP_PORT TCP_ROUTER_DNS_NAME CF_APPS_DOMAIN

  ginkgo "$GOPATH/src/integration-tests/cloudfoundry"
elif [[ ${routing_mode} == "iaas" && ${iaas} == "gcp" ]]; then
  ginkgo "$GOPATH/src/integration-tests/pod_logs"
  ginkgo "$GOPATH/src/integration-tests/workload/k8s_lbs"
elif [[ ${routing_mode} == "iaas" ]]; then
  WORKLOAD_ADDRESS=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_worker_host")
  export WORKLOAD_ADDRESS
  ginkgo "$GOPATH/src/integration-tests/pod_logs"
  ginkgo "$GOPATH/src/integration-tests/workload/iaas_lbs"
elif [[ ${routing_mode} == "proxy" ]]; then
  WORKLOAD_ADDRESS=$(call_bosh -d "${DEPLOYMENT_NAME}" vms | grep 'worker-haproxy/' | head -1 | awk '{print $4}')
  WORKLOAD_PORT=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/worker_haproxy_tcp_frontend_port")
  export WORKLOAD_ADDRESS WORKLOAD_PORT

  ginkgo "$GOPATH/src/integration-tests/workload/haproxy"
fi
ginkgo "$GOPATH/src/integration-tests/generic"
