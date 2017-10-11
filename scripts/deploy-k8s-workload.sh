#!/bin/bash

set -exo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/lb-info.sh"

if [ -z ${LOCAL_DEV+x} ] && [ "$LOCAL_DEV" != "1" ]; then
  cp "gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
fi

bosh_ca_cert=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/default_ca/ca)
client_secret=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/bosh_admin_client_secret)

director_ip=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip")

"$KUBO_DEPLOYMENT_DIR/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service

kubectl create -f "$KUBO_CI_DIR/specs/nginx-lb.yml"
kubectl rollout status -w deployment/nginx

# get the load balancer's address
lb_address_blocking nginx "$KUBO_ENVIRONMENT_DIR" "$KUBO_DEPLOYMENT_DIR"
if [ "$?" != 0 ]; then exit 1; fi

lb_url="http://$lb_address"

timeout_seconds=10
max_attempts=30
current_attempt=0
retry=true

# Probe the workload to ensure that it is servicing requests
while $retry; do
  curl -L --max-time ${timeout_seconds} -IfsS ${lb_url}

  if [ $? -eq 0 ]; then
    retry=false
  else
    current_attempt=$((current_attempt+1))
    if [ $current_attempt -gt $max_attempts ]; then
      echo "Reached maximum attempts trying to query $lb_url"
      exit 1
    fi
  fi
done
