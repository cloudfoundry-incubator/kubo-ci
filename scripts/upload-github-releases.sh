#!/usr/bin/env bash

set -exu -o pipefail

version=$(cat kubo-version/version)
echo "kubo-release ${version}" >kubo-release/name
echo "See [CFCR Release notes](https://docs-cfcr.cfapps.io/overview/release-notes/) page" > kubo-release/body

echo "kubo-deployment ${version}" >kubo-deployment/name
echo "See [CFCR Release notes](https://docs-cfcr.cfapps.io/overview/release-notes/) page" > kubo-deployment/body

mkdir "kubo-deployment-${version}"
cp -r git-kubo-deployment-output/* "kubo-deployment-${version}/"

tar -czf kubo-deployment/kubo-deployment-${version}.tgz "kubo-deployment-${version}"
