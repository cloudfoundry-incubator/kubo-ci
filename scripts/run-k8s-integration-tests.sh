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

function credHub_login {
    local director_name credhub_user_password credhub_api_url
    director_name=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name")
    credhub_user_password=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_cli_password")
    credhub_api_url="https://$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip"):8844"

    tmp_uaa_ca_file="$(mktemp)"
    bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/uaa_ssl/ca" > "${tmp_uaa_ca_file}"
    tmp_credhub_ca_file="$(mktemp)"
    bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_tls/ca" > "${tmp_credhub_ca_file}"

    credhub login -u credhub-cli -p "${credhub_user_password}" -s "${credhub_api_url}" --ca-cert "${tmp_credhub_ca_file}" --ca-cert "${tmp_uaa_ca_file}"
}

GIT_KUBO_DEPLOYMENT_DIR=$1
DEPLOYMENT_NAME=$2
KUBO_ENVIRONMENT_DIR=$3

credHub_login

if [ -z "${SKIP_KUBECONFIG+1}" ]; then
  "$GIT_KUBO_DEPLOYMENT_DIR/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}"
fi

routing_mode=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing_mode")
iaas=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/iaas")
director_name=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name")
GIT_KUBO_CI=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)
GOPATH="$GIT_KUBO_CI"
INTEGRATIONTEST_IAAS=${iaas}

export GOPATH INTEGRATIONTEST_IAAS DEPLOYMENT_NAME

export PATH_TO_KUBECONFIG="$HOME/.kube/config"
TLS_KUBERNETES_CERT=$(bosh-cli int <(credhub get -n "${director_name}/${DEPLOYMENT_NAME}/tls-kubernetes" --output-json) --path='/value/certificate')
TLS_KUBERNETES_PRIVATE_KEY=$(bosh-cli int <(credhub get -n "${director_name}/${DEPLOYMENT_NAME}/tls-kubernetes" --output-json) --path='/value/private_key')
export TLS_KUBERNETES_CERT TLS_KUBERNETES_PRIVATE_KEY

BOSH_ENVIRONMENT=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path='/internal_ip')
BOSH_CA_CERT=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path='/default_ca/ca')
BOSH_CLIENT=bosh_admin
BOSH_CLIENT_SECRET=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path='/bosh_admin_client_secret')

export BOSH_ENVIRONMENT BOSH_CA_CERT BOSH_CLIENT BOSH_CLIENT_SECRET

if [[ ${routing_mode} == "cf" ]]; then
  KUBERNETES_SERVICE_HOST=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_host")
  KUBERNETES_SERVICE_PORT=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_port")
  WORKLOAD_TCP_PORT=$(expr "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_port")" + 10)
  INGRESS_CONTROLLER_TCP_PORT=$(expr "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_port")" + 20)
  TCP_ROUTER_DNS_NAME=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/kubernetes_master_host")
  CF_APPS_DOMAIN=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing-cf-app-domain-name")
  export KUBERNETES_SERVICE_HOST KUBERNETES_SERVICE_PORT WORKLOAD_TCP_PORT INGRESS_CONTROLLER_TCP_PORT TCP_ROUTER_DNS_NAME CF_APPS_DOMAIN

  ginkgo "$GOPATH/src/tests/integration-tests/cloudfoundry"
elif [[ ${routing_mode} == "iaas" ]]; then

  case "${iaas}" in
    aws)
      aws configure set aws_access_key_id "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/access_key_id)"
      aws configure set aws_secret_access_key  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/secret_access_key)"
      aws configure set default.region "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/region)"
      AWS_INGRESS_GROUP_ID=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path=/default_security_groups/0)
      export AWS_INGRESS_GROUP_ID
      ;;
  esac

  ginkgo "$GOPATH/src/tests/integration-tests/workload/k8s_lbs"
elif [[ ${routing_mode} == "proxy" ]]; then
  WORKLOAD_ADDRESS=$(call_bosh -d "${DEPLOYMENT_NAME}" vms | grep 'worker-haproxy/' | head -1 | awk '{print $4}')
  WORKLOAD_PORT=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/worker_haproxy_tcp_frontend_port")
  export WORKLOAD_ADDRESS WORKLOAD_PORT

  ginkgo "$GOPATH/src/tests/integration-tests/workload/haproxy"
fi
ginkgo "$GOPATH/src/tests/integration-tests/pod_logs"
ginkgo "$GOPATH/src/tests/integration-tests/generic"
ginkgo "$GOPATH/src/tests/integration-tests/oss_only"

if [[ "${iaas}" != "openstack" ]]; then
    ginkgo "$GOPATH/src/tests/integration-tests/persistent_volume"
fi
