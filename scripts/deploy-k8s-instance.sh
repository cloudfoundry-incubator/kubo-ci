#!/bin/sh -e

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

tarball_name=$(ls $PWD/s3-kubo-release-tarball/kubo-release*.tgz | head -n1)

cp "$PWD/s3-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"

cp "$tarball_name" "git-kubo-deployment/../kubo-release.tgz"

credhub login -u credhub-user -p \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path="/credhub_user_password")" \
  -s "https://$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip"):8844" --skip-tls-validation
credhub set -n \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name")/ci-service/routing-cf-nats-password" \
  -t password -v "${ROUTING_CF_NATS_PASSWORD}" -O > /dev/null

credhub set -n \
  "$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/director_name")/ci-service/routing-cf-client-secret" \
  -t password -v "${ROUTING_CF_CLIENT_SECRET}" -O > /dev/null


"git-kubo-deployment/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"
# Deploy k8s
"git-kubo-deployment/bin/deploy_k8s" "${KUBO_ENVIRONMENT_DIR}" ci-service local

cp "${KUBO_ENVIRONMENT_DIR}/ci-service-creds.yml" "$PWD/service-creds/ci-service-creds.yml"
