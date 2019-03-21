#!/bin/bash -exu

bosh -d concourse-worker \
  deploy \
  --vars-file=vsphere-worker-vars.yml \
  --vars-file=$HOME/workspace/concourse-bosh-deployment/versions.yml \
  -v external_worker_network_name=kubo-network \
  -v deployment_name=concourse-worker \
  -v azs=[z1] \
  -v instances=2 \
  -v worker_tags=[vsphere-lb,vsphere-proxy] \
  -v worker_vm_type=worker \
  -v tsa_host=ci.kubo.sh \
  $HOME/workspace/concourse-bosh-deployment/cluster/external-worker.yml
