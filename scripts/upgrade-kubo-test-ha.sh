#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"

HA_MIN_SERVICE_AVAILABILITY="${HA_MIN_SERVICE_AVAILABILITY:-1}"

update_kubo() {
  # Workaround due to https://www.pivotaltracker.com/story/show/152155545
  echo "Deleting tls-kubernetes from credhub..."
  credhub api "https://$(bosh-cli int environment/creds.yml --path=/internal_ip):8844"
  credhub login -u credhub-cli -p "$(bosh-cli int environment/creds.yml --path=/credhub_cli_password)"
  credhub delete "$(bosh-cli int environment/director.yml --path=/director_name)/ci-service/tls-kubernetes"

  echo "Updating Kubo..."
  ${DIR}/deploy-k8s-instance.sh
}

upload_new_releases
run_upgrade_test update_kubo "$HA_MIN_SERVICE_AVAILABILITY"
