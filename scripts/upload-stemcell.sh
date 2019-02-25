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


stemcell_version="$(bosh int --path=/stemcells/0/version git-kubo-deployment/manifests/cfcr.yml)"
stemcell_line="$(bosh int --path=/stemcells/0/os git-kubo-deployment/manifests/cfcr.yml)"
bosh upload-stemcell --name="bosh-${IAAS}-${VM}-${stemcell_line}-go_agent" --version="${stemcell_version}" "https://s3.amazonaws.com/bosh-core-stemcells/${IAAS}/bosh-stemcell-${stemcell_version}-${IAAS}-${VM}-${stemcell_line}-go_agent.tgz"

stemcell_version="$(bosh int -o git-kubo-deployment/manifests/ops-files/windows/change-windows-stemcell.yml --path=/stemcells/1/version git-kubo-deployment/manifests/cfcr.yml)"
stemcell_line="$(bosh int -o git-kubo-deployment/manifests/ops-files/windows/change-windows-stemcell.yml --path=/stemcells/1/os git-kubo-deployment/manifests/cfcr.yml)"
bosh upload-stemcell --name="bosh-${IAAS}-${VM}-${stemcell_line}-go_agent" --version="${stemcell_version}" "https://s3.amazonaws.com/bosh-core-stemcells/${IAAS}/bosh-stemcell-${stemcell_version}-${IAAS}-${VM}-${stemcell_line}-go_agent.tgz"