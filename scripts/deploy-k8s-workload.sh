#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

deployment_name="ci-service"

cp "$PWD/gcs-service-creds/ci-service-creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

credhub login -u credhub-user -p \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_user_password")" \
  -s "https://$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip"):8844" --skip-tls-validation

"git-kubo-deployment/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${deployment_name}"

routing_mode=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing_mode")
director_name=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name")

export GOPATH="$PWD/git-kubo-ci"
export PATH_TO_KUBECONFIG="$HOME/.kube/config"

export TLS_KUBERNETES_CERT=$(credhub get -n '${director_name}/${deployment_name}/tls-kubernetes' --output-json | bosh-cli int - --path='/certificate')
export TLS_KUBERNETES_PRIVATE_KEY=$(credhub get -n '${director_name}/${deployment_name}/tls-kubernetes' --output-json | bosh-cli int - --path='/private_key')

if [[ ${routing_mode} == "cf" ]]; then
  export KUBERNETES_SERVICE_HOST=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/cf-tcp-router-name")
  export KUBERNETES_SERVICE_PORT=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/external-kubo-port")
  export WORKLOAD_TCP_PORT=$(expr $(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/external-kubo-port") + 1000)
  export INGRESS_CONTROLLER_TCP_PORT=$(expr $(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/external-kubo-port") + 2000)
  export TCP_ROUTER_DNS_NAME=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/cf-tcp-router-name")
  export CF_APPS_DOMAIN=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing-cf-app-domain-name")

  ginkgo "$GOPATH/src/integration-tests/cloudfoundry"
elif [[ ${routing_mode} == "iaas" ]]; then
  export WORKER_LB_ADDRESS=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_worker_ip")
  export NODE_PORT=$(bosh-cli int "${PWD}/git-kubo-ci/specs/nginx.yml" --path="/spec/ports/0/nodePort")

  ginkgo "$GOPATH/src/integration-tests/gcp_lb"
fi
