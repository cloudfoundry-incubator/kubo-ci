#!/bin/bash
set -eux

cp gcs-shipable-version/* gcs-shipable-version-output
echo | cat kubo-version/number - >> gcs-shipable-version-output/shipable
exit 1
