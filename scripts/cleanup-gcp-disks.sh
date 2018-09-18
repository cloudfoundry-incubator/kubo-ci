#!/bin/bash

set -eux

export GOOGLE_APPLICATION_CREDENTIALS="$(mktemp)"

set +x
echo "$GCP_SERVICE_ACCOUNT" > "$GOOGLE_APPLICATION_CREDENTIALS"
set -x

gcloud auth activate-service-account --key-file "$GOOGLE_APPLICATION_CREDENTIALS" --project cf-pcf-kubo

disks_to_delete="$(gcloud compute disks list --format='value(selfLink)' --filter='-users:*')"

if [[ "$disks_to_delete" != "" ]]; then
    gcloud -q compute disks delete $disks_to_delete
fi
