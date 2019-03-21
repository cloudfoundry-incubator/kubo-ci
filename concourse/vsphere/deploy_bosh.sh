#!/bin/bash -exu

my_dir=$PWD

pushd ~/workspace/bosh-deployment
  bosh create-env bosh.yml \
      --state=$my_dir/state.json \
      --vars-store=$my_dir/creds.yml \
      --vars-file=$my_dir/director.yml \
      -o vsphere/cpi.yml
popd
