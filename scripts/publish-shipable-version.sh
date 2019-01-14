#!/bin/bash
set -exu -o pipefail

signal="$(basename $(find gcs-shipable-version ! -name 'url' ! -name 'generation' -type f))"
cp gcs-shipable-version/$signal gcs-shipable-version-output/shipable
mkdir kubo-release-untarred
tar -xzf kubo-release/kubo-*.tgz --directory kubo-release-untarred
release_hash="$(grep "commit_hash" kubo-release-untarred/release.MF | awk -F ' ' '{print $2}')"
deployment_hash="$(cat git-kubo-deployment/.git/short_ref)"

printf "%s %s\n" "${release_hash}" "${deployment_hash}" >> "gcs-shipable-version-output/shipable"
