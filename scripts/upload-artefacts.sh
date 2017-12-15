#!/bin/bash

set -exu -o pipefail

tar -zxvf gcs-kubo-deployment-tarball-untested/*.tgz
./kubo-deployment/bin/upload_artefacts "kubo-lock" "local"
