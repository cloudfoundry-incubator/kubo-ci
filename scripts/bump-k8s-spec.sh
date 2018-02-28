#!/bin/bash

set -eu -o pipefail

source /common.sh
start_docker

set -x
if [ ! -f spec-to-update/spec.env ]; then
    echo "No new versions found to update."
    exit 0
fi
source spec-to-update/spec.env

tag=$(cat "$PWD/$SPEC_RELEASE_DIR/tag")
version=$(cat "$PWD/$SPEC_RELEASE_DIR/version")

cp -r git-kubo-release/. git-kubo-release-output
pushd git-kubo-release-output

scripts/download_container_images "$SPEC_IMAGE_URL:$tag"
sed -E -i -e "/${SPEC_IMAGE_NAME}:/s/v([0-9]+\.)+[0-9]+/${tag}/" scripts/download_container_images
sed -E -i -e "/${SPEC_IMAGE_NAME}:/s/v([0-9]+\.)+[0-9]+/${tag}/" "jobs/apply-specs/templates/specs/${SPEC_NAME}.yml.erb"

set +x
cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF
set -x

bosh upload-blobs

git config --global user.email "cfcr+cibot@pivotal.io"
git config --global user.name "CFCR CI BOT"
git add .
git commit -m "Bump ${SPEC_NAME} to version ${version}"
popd
