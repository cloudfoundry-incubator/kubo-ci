#!/bin/sh -e

bosh-cli create-env "git-kubo-deployment/bosh-deployment/bosh.yml"  \
  --ops-file "git-kubo-deployment/bosh-deployment/aws/cpi.yml" \
  --ops-file "git-kubo-deployment/bosh-deployment/powerdns.yml" \
  --ops-file "git-kubo-deployment/bosh-deployment/jumpbox-user.yml" \
  --ops-file "git-kubo-ci/etcd/bosh_admin_user_ops_file.yml" \
  --ops-file "git-kubo-ci/etcd/increase_bosh_workers.yml" \
  --state "bosh-state/state.json" \
  --vars-store "bosh-creds/creds.yml" \
  --vars-file "kubo-lock/metadata"
