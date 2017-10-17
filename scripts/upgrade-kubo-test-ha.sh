#!/bin/bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"

HA_MIN_SERVICE_AVAILABILITY="${HA_MIN_SERVICE_AVAILABILITY:-1}"

update_kubo() {
  echo "Updating Kubo..."
  ${DIR}/deploy-k8s-instance.sh
}

upload_new_releases
run_upgrade_test update_kubo "$HA_MIN_SERVICE_AVAILABILITY"
