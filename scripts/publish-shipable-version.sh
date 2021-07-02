#!/bin/bash
set -exu -o pipefail

signal="$(basename $(find gcs-shipable-version ! -name 'url' ! -name 'generation' -type f))"
cp gcs-shipable-version/$signal gcs-shipable-version-output/shipable
mkdir kubo-release-untarred
tar -xzf kubo-release/kubo-*.tgz --directory kubo-release-untarred
release_hash=$(grep "commit_hash" kubo-release-untarred/release.MF | awk -F ' ' '{print $2}' | tr -d "'")
dev_version="$(grep "^version" kubo-release-untarred/release.MF | awk -F ' ' '{print $2}')"

printf "%s %s\n" "${release_hash}" "${dev_version}">> "gcs-shipable-version-output/shipable"
