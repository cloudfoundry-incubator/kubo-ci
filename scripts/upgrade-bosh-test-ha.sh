#!/bin/bash

set -o pipefail

update() {
  echo "Updating BOSH..."
  DO_UPGRADE=1 ./install-bosh.sh
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
  update &
  update_pid="$!"

  url="example.com"
  query_loop "$update_pid" "$url" &
  query_pid="$!"

  wait_for_success "$update_pid" "Update BOSH"
  wait_for_success "$query_pid" "HA query loop"
}

main $@