
#!/bin/bash

set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

. "$DIR/lib/environment.sh"
. "$DIR/lib/upgrade-tests.sh"

update_kubo() {
  echo "Updating Kubo..."
  ${DIR}/deploy-k8s-instance.sh
}

run_upgrade_test update_kubo
