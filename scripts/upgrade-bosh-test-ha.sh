#!/bin/bash

set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "${DIR}/lib/environment.sh"

# copy state and creds so that deploy_bosh has the correct context
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}"
cp "$PWD/gcs-bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
touch "${KUBO_ENVIRONMENT_DIR}/director-secrets.yml"

update() {
  echo "Updating BOSH..."
  ${DIR}/install-bosh.sh
}

query_loop() {
  timeout_seconds=10
  pid_to_wait="$1"
  url="$2"

  echo "Querying $url while waiting for process with pid $pid_to_wait to finish..."

  # loop while process is not finished
  while kill -0 "$pid_to_wait" >/dev/null 2>&1; do
    sleep 1

    response_code=$(curl -L --max-time "$timeout_seconds" -s -o /dev/null -I -w "%{http_code}" "$url")

    if [ "$response_code" != "200" ]; then
      echo "Error: response from $url is not 200 (got $response_code)"
      exit 1
    else
      echo "Service successfully returned response code 200"
    fi
  done
}

# query the API for the load balancer address and sets $lb_address
lb_address_blocking() {
  service_name="$1"

  iaas=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/iaas")

	lb_address=""
  current_attempt=0
  max_attempts=10

  while [ -z "$lb_address" ]; do
    current_attempt=$((current_attempt+1))
    if [ ${current_attempt} -gt ${max_attempts} ]; then
      echo "Error: reached max attempts trying to obtain load balancer IP"
      exit 1
    fi

    "${KUBO_DEPLOYMENT_DIR}/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service

    if [ ${iaas} = "gcp" ]; then
			lb_address=$(kubectl get service ${service_name} -o jsonpath={.status.loadBalancer.ingress[0].ip})
    else
      lb_address=$(kubectl get service ${service_name} -o jsonpath={.status.loadBalancer.ingress[0].hostname})
    fi

    if [ -z "$lb_address" ]; then sleep 10; fi
  done
}

wait_for_success() {
  pid_to_wait="$1"
  work_description="$2"

  wait "$pid_to_wait"
  if [ "$?" -ne 0 ]; then
    echo "$work_description failed"
    exit 1
  fi

  echo "$work_description succeeded"
}

main() {
  service_name="nginx"

  lb_address_blocking ${service_name}
  if [ -z "$lb_address" ]; then
    echo "Error: couldn't obtain load balancer address"
    exit 1
  fi

  lb_url="http://$lb_address"
  echo "The load balancer's URL is $lb_url"

  # update BOSH in the background
  update &
  update_pid="$!"

  # exercise the load balancer URL while BOSH is updating
  query_url="$lb_url"
  query_loop "$update_pid" "$query_url" &
  query_pid="$!"

  wait_for_success "$update_pid" "Update BOSH"
  wait_for_success "$query_pid" "HA query loop"
}

main $@