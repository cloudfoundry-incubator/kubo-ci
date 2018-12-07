#!/bin/bash

set -exu -o pipefail

main() {
    ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}/../..")" && pwd)
    source "${ROOT}/git-kubo-ci/scripts/set-bosh-env" "${VARFILE}"
    deployment_name="$(bosh int "${VARFILE}" --path "/deployment_name")"

    set +x
    bosh -d "${deployment_name}" -n delete-deployment --force

    ${ROOT}/git-kubo-ci/scripts/credhub-login "${VARFILE}"

    # don't delete leading & trailing slash. This is to scope to the deployment creds we want to delete
    credhub find -n "/${DEPLOYMENT_NAME}/" --output-json | jq -r .credentials[].name | xargs -L 1 credhub delete -n
}

main
