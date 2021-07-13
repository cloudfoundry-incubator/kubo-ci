#!/usr/bin/env bash

set -eu
if [[ -f source-json/source.json ]]; then
    source git-kubo-ci/scripts/set-bosh-env source-json/source.json
else
    source git-kubo-ci/scripts/set-bosh-env source-json/metadata
fi

VM=""
if [ ${IAAS} == "gcp" ]; then
  IAAS="google"
  VM="kvm"
elif [ ${IAAS} == "aws" ]; then
  VM="xen-hvm"
elif [ ${IAAS} == "vsphere" ]; then
  VM="esxi"
elif [ ${IAAS} == "azure" ]; then
  VM="hyperv"
elif [ ${IAAS} == "openstack" ]; then
  VM="kvm"
fi


stemcell_version="$(bosh int --path=/stemcells/0/version git-kubo-release/manifests/cfcr.yml)"
stemcell_line="$(bosh int --path=/stemcells/0/os git-kubo-release/manifests/cfcr.yml)"

bosh upload-stemcell --name="bosh-${IAAS}-${VM}-${stemcell_line}-go_agent" \
  --version="${stemcell_version}" \
  "https://storage.googleapis.com/bosh-${IAAS}-light-stemcells/${stemcell_version}/light-bosh-stemcell-${stemcell_version}-${IAAS}-${VM}-${stemcell_line}-go_agent.tgz"

if [[ -d alternate-stemcell ]]; then
  files=( alternate-stemcell/*stemcell*.tgz )
  bosh upload-stemcell "${files[0]}"
fi
