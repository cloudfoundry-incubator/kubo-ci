#!/bin/bash
set -exu -o pipefail

signal="$(basename $(find gcs-shipable-version ! -name 'url' ! -name 'generation' -type f))"
cp gcs-shipable-version/$signal gcs-shipable-version-output/shipable
echo | cat kubo-version/number - >> gcs-shipable-version-output/shipable
