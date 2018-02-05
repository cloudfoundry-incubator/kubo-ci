#!/bin/bash

set -exu -o pipefail

tag=$(cat "$PWD/govc-release/tag")

cp -r git-kubo-ci/. git-kubo-ci-output
pushd git-kubo-ci-output/docker-images/kubo-ci
sed -E -i -e "/govmomi/s/v([0-9]+\.)+[0-9]+/${tag}/" Dockerfile

git config --global user.email "cfcr+cibot@pivotal.io"
git config --global user.name "CFCR CI BOT"
git add .
git commit -m "Bump CI image to govc-$tag"
popd
