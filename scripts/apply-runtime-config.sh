#!/bin/bash -exu

SPEC='
---
name: etc_hosts

templates:
  pre_start.erb: bin/pre-start

packages: []

properties:
  etc_hosts:
    description: "/etc/hosts entries"
'

PRE_START_ERB='#!/bin/bash -ex

<% p("etc_hosts").each do |entry| %>
  echo "<%= entry %>" >> /etc/hosts
<% end %>
'

FINAL_YML='
name: etc_hosts
'

BLOBS_YML='
--- {}
'

RUNTIME_CONFIG_YML='
releases:
- {name: etc_hosts, version: 42}

addons:
- name: etc_hosts
  jobs:
  - name: etc_hosts
    release: etc_hosts
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

upload_etc_hosts_release() {
  pushd "$(mktemp -d)"
    mkdir -p config packages src
    mkdir -p jobs/etc_hosts/templates

    echo "$SPEC" > jobs/etc_hosts/spec
    echo "$PRE_START_ERB" > jobs/etc_hosts/templates/pre_start.erb
    echo "$FINAL_YML" > config/final.yml
    echo "$BLOBS_YML" > config/blobs.yml
    touch jobs/etc_hosts/monit

    bosh-cli create-release --version 42 && bosh-cli upload-release
  popd
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

  upload_etc_hosts_release

  update_runtime_config
}

main "$@"
