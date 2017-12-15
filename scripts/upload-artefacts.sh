#!/bin/bash

set -exu -o pipefail

mkdir kubo-deployment
tar -zxvf gcs-kubo-deployment-tarball-untested/*.tgz -C kubo-deployment
./kubo-deployment/bin/upload_artefacts "kubo-lock" "local"