#!/usr/bin/env bash

set -exu -o pipefail

main() {
  kubo_release_version=$(cat kubo-release/version)
  kubo_deployment_version=$(cat kubo-deployment/version)
  if [ "$kubo_release_version" != "$kubo_deployment_version" ]; then
    echo "kubo-release version: $kubo_release_version doesn't match kubo-deployment version: $kubo_deployment_version"
    exit 1
  fi
}

main @
