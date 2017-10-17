#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"

HA_MIN_SERVICE_AVAILABILITY="${HA_MIN_SERVICE_AVAILABILITY:-1}"

update_stemcell() {
  echo "Updating Stemcell..."

  echo "Updating kubo-deployment manifest..."
  local new_stemcell_version="3445.11"
  local manifest_path="${KUBO_DEPLOYMENT_DIR}/manifests/kubo.yml"
  ruby -e "require 'yaml'; data = YAML.load_file(\"${manifest_path}\"); data[\"stemcells\"][0][\"version\"] = \"${new_stemcell_version}\"; File.open(\"${manifest_path}\", 'w') { |f| f.write(data.to_yaml.gsub(\"---\n\", \"\")) }"

  ${DIR}/deploy-k8s-instance.sh
}

upload_new_releases
run_upgrade_test update_stemcell "$HA_MIN_SERVICE_AVAILABILITY"
