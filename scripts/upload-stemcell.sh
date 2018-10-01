#!/usr/bin/env bash

set -eu

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" && pwd)
KUBO_ENVIRONMENT_DIR="${ROOT}/environment"

mkdir -p "${KUBO_ENVIRONMENT_DIR}"
cp "${ROOT}/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "${ROOT}/kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

BOSH_ENV="${KUBO_ENVIRONMENT_DIR}" source "${ROOT}/git-kubo-ci/scripts/set_bosh_environment"
stemcell_version="$(bosh int --path=/stemcells/0/version $ROOT/git-kubo-deployment/manifests/cfcr.yml)"
stemcell_line="$(bosh int --path=/stemcells/0/os $ROOT/git-kubo-deployment/manifests/cfcr.yml)"

bosh upload-stemcell "https://boshstemcells.com/${IAAS}/${stemcell_line}/${stemcell_version}"
