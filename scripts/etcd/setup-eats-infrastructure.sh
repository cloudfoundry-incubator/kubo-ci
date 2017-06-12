#!/bin/sh -e

set +x
bosh-cli int kubo-lock/metadata --path=/gcp_service_account > "key.json"
set -x

bosh-cli create-env "git-kubo-deployment/bosh-deployment/bosh.yml"  \
  --ops-file "git-kubo-deployment/bosh-deployment/gcp/cpi.yml" \
  --ops-file "git-kubo-deployment/bosh-deployment/powerdns.yml" \
  --ops-file "git-kubo-ci/etcd/bosh_admin_user_ops_file.yml" \
  --state "bosh-state/state.json" \
  --var-file gcp_credentials_json=key.json \
  --vars-store "bosh-creds/creds.yml" \
  --vars-file "kubo-lock/metadata"
