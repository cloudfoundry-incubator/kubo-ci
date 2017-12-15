#!/bin/bash

set -exo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/lb-info.sh"

if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
fi

bosh_ca_cert=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/default_ca/ca)
client_secret=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/bosh_admin_client_secret)

director_ip=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip")

"$KUBO_DEPLOYMENT_DIR/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service


randomString() {
  head /dev/urandom | tr -dc A-Za-z0-9 | head -c 13 ; echo ''
}

# get the load balancer's address
routing_mode=$(bosh-cli int environment/director.yml --path=/routing_mode)
if  [[ "$routing_mode" == "iaas" ]]; then
  kubectl apply -f "$KUBO_CI_DIR/specs/nginx-lb.yml"
  kubectl rollout status -w deployment/nginx
  lb_address_blocking nginx "$KUBO_ENVIRONMENT_DIR" "$KUBO_DEPLOYMENT_DIR"
  if [ "$?" != 0 ]; then exit 1; fi
elif [[ "$routing_mode" == "cf" ]]; then
  kubectl apply -f "$KUBO_CI_DIR/specs/nginx.yml"
  kubectl rollout status -w deployment/nginx
  service_name=$(randomString)
  kubectl label services nginx "http-route-sync=$service_name" --overwrite
  cf_apps_domain=$(bosh-cli int environment/director.yml --path=/routing_cf_app_domain_name)
  lb_address="$service_name"."$cf_apps_domain"
else
  echo "Routing mode '$routing_mode' is not supported in this test"
  exit 1
fi

lb_url="http://$lb_address"

timeout_seconds=10
max_attempts=30
current_attempt=0
retry=true

# Probe the workload to ensure that it eventually comes up
while $retry; do
  if curl -L --max-time ${timeout_seconds} -IfsS ${lb_url}; then
    retry=false
  else
    current_attempt=$((current_attempt+1))
    if [ $current_attempt -gt $max_attempts ]; then
      echo "Reached maximum attempts trying to query $lb_url"
      exit 1
    fi
    sleep 1
  fi
done
