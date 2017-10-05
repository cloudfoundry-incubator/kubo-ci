#!/bin/bash

set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

update() {
  echo "Updating BOSH..."
  DO_UPGRADE=1 ${DIR}/install-bosh.sh
}

query_loop() {
  timeout=10
  pid_to_wait="$1"
  url="$2"

  echo "Querying $url while waiting for process with pid $pid_to_wait to finish..."

  while kill -0 "$pid_to_wait" >/dev/null 2>&1; do
    sleep 1

    response_code=$(curl -L --max-time "$timeout" -s -o /dev/null -I -w "%{http_code}" "$url")

    if [ "$response_code" -ne 200 ]; then
      echo "Error: response from $url is not 200 (got $response_code)"
      exit 1
    fi
  done
}

lb_ip_blocking() {
  service_name="$1"
	lb_ip=""
  current_attempt=0
  max_attempts=10
  
  while [ -z "$lb_ip" ]; do
    current_attempt=$((current_attempt+1))
    if [ ${current_attempt} -gt ${max_attempts} ]; then
      echo "Error: reached max attempts trying to obtain load balancer IP"
      exit 1
    fi

    # AWS specific?
    lb_ip=$(kubectl get service ${service_name} -o jsonpath={.status.loadBalancer.ingress[0].hostname})

    if [ -z "$lb_ip" ]; then sleep 10; fi
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

  lb_ip_blocking ${service_name}
  if [ -z "$lb_ip" ]; then
    echo "Error: couldn't obtain load balancer IP"
    exit 1
  fi

  lb_url="http://$lb_ip"
  echo "The load balancer's URL is $lb_url"

  update &
  update_pid="$!"

  query_url="$lb_url"
  query_loop "$update_pid" "$query_url" &
  query_pid="$!"

  wait_for_success "$update_pid" "Update BOSH"
  wait_for_success "$query_pid" "HA query loop"
}

main $@