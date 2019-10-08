#!/bin/bash -exu

mkdir -p /tmp/stemcell/
wget -O /tmp/stemcell/stemcell.tgz https://s3.amazonaws.com/bosh-gce-light-stemcells/456.16/light-bosh-stemcell-456.16-google-kvm-ubuntu-xenial-go_agent.tgz
echo "456.16" > /tmp/stemcell/version

mkdir -p /tmp/git-kubo-deployment-output

piper -c tasks/bump-stemcell.yml \
    -i git-kubo-ci=/Users/fulton/workspace/kubo-ci \
    -i git-kubo-deployment=/Users/fulton/workspace/kubo-deployment \
    -i stemcell=/tmp/stemcell \
    -o git-kubo-deployment-output=/tmp/git-kubo-deployment-output