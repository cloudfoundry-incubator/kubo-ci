#!/bin/bash

set -exu -o pipefail

main() {
    ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
    source "${ROOT}/git-kubo-ci/scripts/set-bosh-env" "${VARFILE}"
    deployment_name="$(bosh int "${VARFILE}" --path "/deployment_name")"
    source ${ROOT}/git-kubo-ci/scripts/credhub-login "${VARFILE}"
    ${ROOT}/git-kubo-ci/scripts/set_kubeconfig_long_lived_director

    # Deployment might be deleted already or broken
    set +e
    kubectl delete ns --all
    kubectl delete pvc --all
    kubectl delete pv --all
    kubectl delete svc --all
    set -e

    set +x
    bosh -d "${deployment_name}" -n delete-deployment --force

    # don't delete leading & trailing slash. This is to scope to the deployment creds we want to delete
    credhub find -n "/${deployment_name}/" --output-json | jq -r .credentials[].name | xargs -L 1 credhub delete -n
}

main
