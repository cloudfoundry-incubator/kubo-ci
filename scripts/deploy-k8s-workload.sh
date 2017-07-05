#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

deployment_name="ci-service"

cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

"git-kubo-deployment/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${deployment_name}"

routing_mode=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing_mode")
director_name=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name")

export GOPATH="$PWD/git-kubo-ci"
export PATH_TO_KUBECONFIG="$HOME/.kube/config"

TLS_KUBERNETES_CERT=$(bosh-cli int <(credhub get -n "${director_name}/${deployment_name}/tls-kubernetes" --output-json) --path='/value/certificate')
TLS_KUBERNETES_PRIVATE_KEY=$(bosh-cli int <(credhub get -n "${director_name}/${deployment_name}/tls-kubernetes" --output-json) --path='/value/private_key')
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
  WORKER_LB_ADDRESS=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_worker_host")
  NODE_PORT=$(bosh-cli int "${PWD}/git-kubo-ci/specs/nginx.yml" --path="/spec/ports/0/nodePort")
  export WORKER_LB_ADDRESS NODE_PORT

  ginkgo "$GOPATH/src/integration-tests/gcp_lb"
fi
