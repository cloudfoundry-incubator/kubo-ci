#!/bin/bash

set -exo pipefail

. git-kubo-ci/scripts/lib/environment.sh

cp "gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

bosh_ca_cert=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/default_ca/ca)
client_secret=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/bosh_admin_client_secret)

director_ip=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip")

"git-kubo-deployment/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service

kubectl create -f git-kubo-deployment/specs/nginx-specified-nodeport.yml
kubectl rollout status -w deployment/nginx