#!/usr/bin/env bash

set -eu

source git-kubo-ci/scripts/set-bosh-env gcs-source-json/source.json
stemcell_version="$(bosh int --path=/stemcells/0/version git-kubo-deployment/manifests/cfcr.yml)"
stemcell_line="$(bosh int --path=/stemcells/0/os git-kubo-deployment/manifests/cfcr.yml)"

bosh upload-stemcell "https://boshstemcells.com/${IAAS}/${stemcell_line}/${stemcell_version}"
