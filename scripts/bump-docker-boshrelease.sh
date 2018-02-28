#!/bin/bash

set -exu -o pipefail

base_dir="$PWD"
tag=$(cat "$base_dir/docker-boshrelease/tag")
version=$(cat "$base_dir/docker-boshrelease/version")
release_sha1=$(sha1sum "$base_dir/docker-boshrelease/docker-${version}.tgz")
download_url="https://github.com/cloudfoundry-incubator/docker-boshrelease/releases/download/${tag}/docker-${version}.tgz"

cp -r git-kubo-deployment/. git-kubo-deployment-output
pushd git-kubo-deployment-output

existing_version=$(bosh int manifests/cfcr.yml --path=/releases/name=docker/version)
sed -E -i -e "s/version: ${existing_version}/version: ${version}/" manifests/cfcr.yml

existing_url=$(bosh int manifests/cfcr.yml --path=/releases/name=docker/url)
sed -E -i -e "s,${existing_url},${download_url},g" manifests/cfcr.yml

exisitng_sha1=$(bosh int manifests/cfcr.yml --path=/releases/name=docker/sha1)
sed -E -i -e "s/${exisitng_sha1}/${release_sha1}/" manifests/cfcr.yml

git config --global user.email "cfcr+cibot@pivotal.io"
git config --global user.name "CFCR CI BOT"
git add .
git commit -m "Bump docker-boshrelease to $tag"

popd
