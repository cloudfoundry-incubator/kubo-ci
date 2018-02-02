#!/bin/bash

set -exu -o pipefail

mkdir tmp
cp -r git-kubo-deployment tmp/kubo-deployment
cp gcs-kubo-release-tarball-untested/kubo-*.tgz tmp/kubo-release.tgz

cd tmp
tar --exclude='src/kubo-deployment-tests' \
    --exclude='.git' \
    --exclude='.gitignore' \
    --exclude='bosh-deployment/.gitignore' \
    --exclude='bosh-deployment/.gitrepo' \
    --exclude='bosh-deployment/tests/.gitignore' \
    -zcvf "../kubo-deployment-tarball/kubo-deployment-$(cat ../kubo-version/version).tgz" *
