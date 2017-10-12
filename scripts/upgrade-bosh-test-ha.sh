#!/bin/bash

set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"

update_bosh() {
  echo "Updating BOSH..."
  ${DIR}/install-bosh.sh
}

run_upgrade_test update_bosh
