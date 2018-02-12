#!/usr/bin/env bash

set -eu

start-bosh

source /tmp/local-bosh/director/env

#
# stemcell metadata/upload
#

STEMCELL_OS=$(bosh int git-kubo-deployment/manifests/cfcr.yml --path /stemcells/0/os)
STEMCELL_VERSION=$(bosh int git-kubo-deployment/manifests/cfcr.yml --path /stemcells/0/version)


bosh -n upload-stemcell "https://s3.amazonaws.com/bosh-core-stemcells/warden/bosh-stemcell-$STEMCELL_VERSION-warden-boshlite-$STEMCELL_OS-go_agent.tgz"

#
# release metadata/upload
#

cd release
tar -xzf *.tgz $( tar -tzf *.tgz | grep 'release.MF' )
RELEASE_NAME=$( grep -E '^name: ' release.MF | awk '{print $2}' | tr -d "\"'" )
RELEASE_VERSION=$( grep -E '^version: ' release.MF | awk '{print $2}' | tr -d "\"'" )

bosh -n upload-release *.tgz
cd ../

#
# compilation deployment
#

cat > manifest.yml <<EOF
---
name: compilation
releases:
- name: "$RELEASE_NAME"
  version: "$RELEASE_VERSION"
stemcells:
- alias: default
  os: "$STEMCELL_OS"
  version: "$STEMCELL_VERSION"
update:
  canaries: 1
  max_in_flight: 1
  canary_watch_time: 1000 - 90000
  update_watch_time: 1000 - 90000
instance_groups: []
EOF

bosh -n -d compilation deploy manifest.yml
bosh -d compilation export-release "$RELEASE_NAME/$RELEASE_VERSION" "$STEMCELL_OS/$STEMCELL_VERSION"

mv *.tgz compiled-release/"$( echo *.tgz | sed "s/\\.tgz$/-$( date -u +%Y%m%d%H%M%S ).tgz/" )"
sha1sum compiled-release/*.tgz
