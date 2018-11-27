#!/usr/bin/env bash

set -eu

if [ -f gcs-source-json/source.json ]; then
    source git-kubo-ci/scripts/set-bosh-env gcs-source-json/source.json
else
    KUBO_ENVIRONMENT_DIR="environment"

    mkdir -p "${KUBO_ENVIRONMENT_DIR}"
    cp "gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
    cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
    BOSH_ENV="${KUBO_ENVIRONMENT_DIR}" source "git-kubo-ci/scripts/set_bosh_environment"
fi
stemcell_version="$(bosh int --path=/stemcells/0/version git-kubo-deployment/manifests/cfcr.yml)"
stemcell_line="$(bosh int --path=/stemcells/0/os git-kubo-deployment/manifests/cfcr.yml)"

bosh upload-stemcell "https://boshstemcells.com/${IAAS}/${stemcell_line}/${stemcell_version}"
