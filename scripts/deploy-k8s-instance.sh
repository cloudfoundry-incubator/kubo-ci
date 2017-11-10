#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1


metadata_path="${KUBO_ENVIRONMENT_DIR}/director.yml"
if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
  cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
  cp "kubo-lock/metadata" "$metadata_path"
  tarball_name=$(ls $PWD/gcs-kubo-release-tarball/kubo-*.tgz | head -n1)
else
  tarball_name="$KUBO_RELEASE_TARBALL"
fi

cp "$tarball_name" "$KUBO_DEPLOYMENT_DIR/../kubo-release.tgz"

"$KUBO_DEPLOYMENT_DIR/bin/set_bosh_alias" "${KUBO_ENVIRONMENT_DIR}"

iaas=$(bosh-cli int $metadata_path --path=/iaas)
iaas_cc_opsfile="$KUBO_CI_DIR/manifests/ops-files/${iaas}-k8s-cloud-config.yml"

CLOUD_CONFIG_OPS_FILE=${CLOUD_CONFIG_OPS_FILE:-""}
if [[ -f "$KUBO_CI_DIR/manifests/ops-files/$CLOUD_CONFIG_OPS_FILE" ]]; then
  CLOUD_CONFIG_OPS_FILES="$KUBO_CI_DIR/manifests/ops-files/$CLOUD_CONFIG_OPS_FILE"
elif [[ -f "$iaas_cc_opsfile" ]]; then
  CLOUD_CONFIG_OPS_FILES="${iaas_cc_opsfile}"
fi
export CLOUD_CONFIG_OPS_FILES

release_source="local"

"$KUBO_DEPLOYMENT_DIR/bin/deploy_k8s" "$KUBO_ENVIRONMENT_DIR" ci-service "$release_source"
