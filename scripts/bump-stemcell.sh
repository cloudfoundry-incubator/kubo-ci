#!/usr/bin/env bash

set -eux -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" && pwd)
cp -r "${ROOT}/git-kubo-deployment/." "${ROOT}/git-kubo-deployment-output"

cat << EOF > replace-stemcell-version.yml
- type: replace
  path: /stemcells/0/version
  value: ((stemcell_version))
EOF
stemcell_version="$(cat stemcell/version)"

bosh int "${ROOT}/git-kubo-deployment/manifests/cfcr.yml" \
  -o replace-stemcell-version.yml \
  -v stemcell_version="${stemcell_version}" \
  > git-kubo-deployment-output/manifests/cfcr.yml

git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"
cd "${ROOT}/git-kubo-deployment-output"

git add .
git commit -m "Bumping stemcell to v${stemcell_version}"
