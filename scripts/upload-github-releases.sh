#!/usr/bin/env bash

set -exu -o pipefail

git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"
version=$(cat kubo-version/version)
echo "kubo-release ${version}" >kubo-release/name
echo "See [CFCR Release notes](https://docs-cfcr.cfapps.io/overview/release-notes/) page" > kubo-release/body

echo "kubo-deployment ${version}" >kubo-deployment/name
echo "See [CFCR Release notes](https://docs-cfcr.cfapps.io/overview/release-notes/) page" > kubo-deployment/body

mkdir "kubo-deployment-${version}"
cp kubo-release-tarball/kubo-release-${version}.tgz kubo-deployment-${version}/kubo-release.tgz
cp -r git-kubo-deployment-output "kubo-deployment-${version}/kubo-deployment"

tar -czf kubo-deployment/kubo-deployment-${version}.tgz "kubo-deployment-${version}"

cd git-kubo-release-master-output
git checkout -b tmp/release
git add .
git commit -m "Final release for v${version}"
git tag -a "v${version}" -m "Tag for version v${version}"
git checkout "${BRANCH:-master}"
git merge tmp/release -m "Merge release branch for v${version}"
git branch -d tmp/release
