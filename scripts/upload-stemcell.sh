#!/usr/bin/env bash

set -eu
if [[ -f gcs-source-json/source.json ]]; then
    source git-kubo-ci/scripts/set-bosh-env gcs-source-json/source.json
else
    source git-kubo-ci/scripts/set-bosh-env gcs-source-json/metadata
    jumpbox_ssh_key="$(mktemp)"
    bosh int --path=/jumpbox_ssh_key gcs-source-json/metadata > ${jumpbox_ssh_key}
    proxy="ssh+socks5://jumpbox@$(bosh int --path=/jumpbox_url gcs-source-json/metadata):22?private-key=${jumpbox_ssh_key}"
    export BOSH_ALL_PROXY="${proxy}"
    export CREDHUB_ALL_PROXY="${proxy}"
fi
stemcell_version="$(bosh int --path=/stemcells/0/version git-kubo-deployment/manifests/cfcr.yml)"
stemcell_line="$(bosh int --path=/stemcells/0/os git-kubo-deployment/manifests/cfcr.yml)"

bosh upload-stemcell "https://boshstemcells.com/${IAAS}/${stemcell_line}/${stemcell_version}"
