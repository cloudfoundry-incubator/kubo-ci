#!/bin/bash

set -exu -o pipefail

mkdir tmp
cp -r git-kubo-deployment tmp/kubo-deployment

cd tmp
tar --exclude='src/kubo-deployment-tests' \
    --exclude='.gitignore' \
    -zcvf "../kubo-deployment-tarball/kubo-deployment-$(cat ../kubo-version/version).tgz" *
