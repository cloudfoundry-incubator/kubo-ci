#!/bin/bash

set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/lb-info.sh"

# copy state and creds so that deploy_bosh has the correct context
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}"
  cp "$PWD/gcs-bosh-state/state.json" "${KUBO_ENVIRONMENT_DIR}"
  cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
  touch "${KUBO_ENVIRONMENT_DIR}/director-secrets.yml"
fi

update() {
  echo "Updating BOSH..."
  ${DIR}/install-bosh.sh
}

query_loop() {
  timeout_seconds=20
  pid_to_wait="$1"
  url="$2"

  echo "Querying $url while waiting for process with pid $pid_to_wait to finish..."

  # loop while process is not finished
  while kill -0 "$pid_to_wait" >/dev/null 2>&1; do
    sleep 1

    curl -L --max-time ${timeout_seconds} -IfsS ${url} &> query_loop_last_output.txt

    if [ "$?" -ne 0 ]; then
      echo "Error: request to $url failed"
      cat query_loop_last_output.txt
      return 1
    else
      echo "Service $url successfully responded"
    fi
  done
}

wait_for_success() {
  pid_to_wait="$1"
  work_description="$2"

  echo "PID to wait on: $pid_to_wait"
  wait "$pid_to_wait"
  if [ "$?" -ne 0 ]; then
    echo "$work_description failed"
    return 1
  fi

  echo "$work_description succeeded"
}

main() {
  service_name="nginx"

  lb_address_blocking "$service_name" "$KUBO_ENVIRONMENT_DIR" "$KUBO_DEPLOYMENT_DIR"
  if [ -z "$lb_address" ]; then
    echo "Error: couldn't obtain load balancer address"
    return 1
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
  if [ "$?" -ne 0 ]; then
    echo "Output of last query below:"
    cat query_loop_last_output.txt
  fi
}

main $@