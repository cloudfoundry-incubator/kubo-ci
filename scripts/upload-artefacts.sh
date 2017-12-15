#!/bin/bash

set -exu -o pipefail

tar -zxvf gcs-kubo-deployment-tarball-untested/*.tgz
cp bosh-creds/creds.yml kubo-lock
mv kubo-lock/metadata kubo-lock/director.yml

./kubo-deployment/bin/upload_artefacts "kubo-lock" "local"
