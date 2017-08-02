#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

if [[ $# -lt 3 ]]; then
    echo "Usage:" >&2
    echo "$0 GIT_KUBO_DEPLOYMENT_DIR DEPLOYMENT_NAME KUBO_ENVIRONMENT_DIR" >&2
    exit 1
fi
# 
# function call_bosh {
#   BOSH_ENV="$KUBO_ENVIRONMENT_DIR" source "$GIT_KUBO_DEPLOYMENT_DIR/bin/set_bosh_environment"
#   bosh-cli "$@"
# }
#
# GIT_KUBO_DEPLOYMENT_DIR=$1
# DEPLOYMENT_NAME=$2
# KUBO_ENVIRONMENT_DIR=$3
#
# "$GIT_KUBO_DEPLOYMENT_DIR/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}"
#
# routing_mode=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/routing_mode")
# director_name=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name")
# GIT_KUBO_CI=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)
# GOPATH="$GIT_KUBO_CI"
# export GOPATH
#
# export PATH_TO_KUBECONFIG="$HOME/.kube/config"
# TLS_KUBERNETES_CERT=$(bosh-cli int <(credhub get -n "${director_name}/${DEPLOYMENT_NAME}/tls-kubernetes" --output-json) --path='/value/certificate')
# TLS_KUBERNETES_PRIVATE_KEY=$(bosh-cli int <(credhub get -n "${director_name}/${DEPLOYMENT_NAME}/tls-kubernetes" --output-json) --path='/value/private_key')
# export TLS_KUBERNETES_CERT TLS_KUBERNETES_PRIVATE_KEY
