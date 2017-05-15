#!/bin/bash -exu


RUNTIME_CONFIG_YML='
releases:
- {name: os-conf, version: 42, url: "https://storage.googleapis.com/kubo-public/os-conf-with-etc-hosts.tgz"}

addons:
- name: etc_hosts
  jobs:
  - name: etc_hosts
    release: os-conf
  properties:
    etc_hosts: ["1.1.1.1 gcr.io", "1.2.3.4 hub.docker.com"]
'

error() {
  echo "$1"
  echo
  exit 1
}

usage() {
  echo "USAGE: $0 [PATH_TO_LOCK_FILE] [PATH_TO_BOSH_CREDS]"
  echo
  exit 1
}

target_bosh_director() {
  export BOSH_ENVIRONMENT=$(bosh-cli int $LOCK_FILE_PATH --path '/internal_ip')
  export BOSH_CLIENT=admin
  export BOSH_CLIENT_SECRET=$(bosh-cli int $CREDS_FILE_PATH --path '/admin_password')
  export BOSH_CA_CERT=$(bosh-cli int $CREDS_FILE_PATH --path '/default_ca/ca')
}

update_runtime_config() {
  bosh-cli -n update-runtime-config <(echo "$RUNTIME_CONFIG_YML")
}

LOCK_FILE_PATH="$1"
CREDS_FILE_PATH="$2"

main() {
  [ ! -f "$LOCK_FILE_PATH" ] && usage
  [ ! -f "$CREDS_FILE_PATH" ] && usage

  target_bosh_director

  update_runtime_config
}

main "$@"
