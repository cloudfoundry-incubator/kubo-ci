#!/bin/bash

set -euxo pipefail

stemcell_version=$(cat ./gcp-stemcell/version)
./git-kubo-deployment/bin/update_stemcell "${stemcell_version}"
git \
  -C ./git-kubo-deployment \
  commit -am "Update stemcell version to $stemcell_version"
