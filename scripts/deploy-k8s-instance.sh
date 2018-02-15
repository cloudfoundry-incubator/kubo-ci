#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

metadata_path="${KUBO_ENVIRONMENT_DIR}/director.yml"
if [[ -z ${LOCAL_DEV+x} ]] || [[ "$LOCAL_DEV" != "1" ]]; then
  cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "kubo-lock/metadata" "$metadata_path"
  tarball_name=$(ls $PWD/gcs-kubo-release-tarball/kubo-*.tgz | head -n1)

  # Copy guestbook if WITHOUT_ADDONS isn't set to true
  if [[ -z ${WITHOUT_ADDONS+x} ]] || [[ "$WITHOUT_ADDONS" != "1" ]]; then
    cp "git-kubo-ci/specs/guestbook.yml" "${KUBO_ENVIRONMENT_DIR}/addons.yml"
  else
    # Delete the addons_spec_path from director.yml
    sed -i.bak '/^addons_spec_path:/d' ${metadata_path}
  fi

else
  tarball_name="$KUBO_RELEASE_TARBALL"
fi

if [[ -z ${WITH_PRIVILEGED_CONTAINERS+x} ]] || [[ "$WITH_PRIVILEGED_CONTAINERS" == "1" ]]; then
  echo "allow_privileged_containers: true" >> "${metadata_path}"
fi

cp "$tarball_name" "$KUBO_DEPLOYMENT_DIR/../kubo-release.tgz"

"$KUBO_DEPLOYMENT_DIR/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"

release_source="local"

DEPLOYMENT_OPS_FILE=${DEPLOYMENT_OPS_FILE:-""}
if [[ -f "${KUBO_CI_DIR}/manifests/ops-files/${DEPLOYMENT_OPS_FILE}" ]]; then
  cp "${KUBO_CI_DIR}/manifests/ops-files/${DEPLOYMENT_OPS_FILE}" "${KUBO_ENVIRONMENT_DIR}/${DEPLOYMENT_NAME}.yml"
fi

"$KUBO_DEPLOYMENT_DIR/bin/deploy_k8s" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}" "$release_source"

"$KUBO_DEPLOYMENT_DIR/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}"
mv ~/.kube/config gcs-kubeconfig/config
