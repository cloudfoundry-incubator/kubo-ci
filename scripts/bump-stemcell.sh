#!/usr/bin/env bash

set -eux -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../" && pwd)
cp -r "${ROOT}/git-kubo-release/." "${ROOT}/git-kubo-release-output"

cat << EOF > replace-stemcell-version.yml
- type: replace
  path: /stemcells/0/version
  value: ((stemcell_version))
EOF
stemcell_version="$(cat stemcell/version)"

bosh int "${ROOT}/git-kubo-release/manifests/cfcr.yml" \
  -o replace-stemcell-version.yml \
  -v stemcell_version="\"${stemcell_version}\"" \
  > git-kubo-release-output/manifests/cfcr.yml

git config --global user.name "cfcr"
git config --global user.email "cfcr@pivotal.io"
cd "${ROOT}/git-kubo-release-output"

if [ -n "$(git status --porcelain)" ]; then
  git add manifests/cfcr.yml
  git commit -m "Bumping stemcell to v${stemcell_version}"
else
  echo "No changes to commit."
fi
