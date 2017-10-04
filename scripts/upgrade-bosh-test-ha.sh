#!/bin/bash

set -eo pipefail

upgrade() {
  echo upgrading...
  sleep 3
}

query() {
  echo querying...
  timeout=10
  url="example.com"
  response_code=$(curl -L --max-time "$timeout" -s -o /dev/null -I -w "%{http_code}" "$url")

  if [ "$response_code" -ne 200 ]; then
    echo "error: response from $url is not 200 (got $response_code)"
    exit 1
  fi
}

main() {
  upgrade &
  upgrade_pid="$!"

  # loop as long as upgrade is ongoing
  while kill -0 "$upgrade_pid" >/dev/null 2>&1; do
    sleep 1
    query
  done
}

main $@