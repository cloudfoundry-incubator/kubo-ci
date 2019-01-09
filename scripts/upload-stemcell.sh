#!/usr/bin/env bash

set -eu
if [[ -f source-json/source.json ]]; then
    source git-kubo-ci/scripts/set-bosh-env source-json/source.json
else
    source git-kubo-ci/scripts/set-bosh-env source-json/metadata
fi

if [ ${IAAS} == "gcp" ]; then
  IAAS="google"
fi

stemcell_version="$(bosh int --path=/stemcells/0/version git-kubo-deployment/manifests/cfcr.yml)"
stemcell_line="$(bosh int --path=/stemcells/0/os git-kubo-deployment/manifests/cfcr.yml)"
bosh upload-stemcell --name="bosh-${IAAS}-kvm-${stemcell_line}-go_agent" --version="${stemcell_version}" "https://s3.amazonaws.com/bosh-core-stemcells/${IAAS}/bosh-stemcell-${stemcell_version}-${IAAS}-xen-hvm-${stemcell_line}-go_agent.tgz"
