#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"

HA_MIN_SERVICE_AVAILABILITY="${HA_MIN_SERVICE_AVAILABILITY:-1}"

if ([ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]) || [ -z "$BOSH_STEMCELL_VERSION" ]; then
  BOSH_STEMCELL_VERSION=$(cat ${PWD}/new-bosh-stemcell/version)
fi

update_stemcell() {
  local manifest_path="${KUBO_DEPLOYMENT_DIR}/manifests/kubo.yml"
  local existing_version=$(bosh int $manifest_path --path=/stemcells/0/version)

  echo "Updating $manifest_path's stemcell version from '$existing_version' to '$BOSH_STEMCELL_VERSION'"
  ruby -e "require 'yaml'; data = YAML.load_file(\"${manifest_path}\"); data[\"stemcells\"][0][\"version\"] = \"${BOSH_STEMCELL_VERSION}\"; File.open(\"${manifest_path}\", 'w') { |f| f.write(data.to_yaml.gsub(\"---\n\", \"\")) }"

  echo "Updating Stemcell..."
  ${DIR}/deploy-k8s-instance.sh
}

upload_new_releases
run_upgrade_test update_stemcell "$HA_MIN_SERVICE_AVAILABILITY"
