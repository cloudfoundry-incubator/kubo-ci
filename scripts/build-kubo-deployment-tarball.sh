#!/bin/bash

set -exu -o pipefail

mkdir tmp
cp -r git-kubo-deployment tmp/kubo-deployment
cp gcs-kubo-release-tarball-untested/kubo-*.tgz tmp/kubo-release.tgz

cd tmp
tar -zcvf "../kubo-deployment-tarball/kubo-deployment-$(cat ../kubo-version/version).tgz" *