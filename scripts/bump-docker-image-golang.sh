#!/bin/bash

set -exu -o pipefail

cp -r git-kubo-ci/. git-kubo-ci-output

pushd "$HOME_DIR"/golang-release
  blob_name=$(bosh blobs --json | jq '.Tables[0].Rows[] | .path | select(test("'"${MINOR_GO_VERSION}"'.*linux"))' --raw-output)
  go_version="${blob_name%.tar.gz}"
popd

pushd git-kubo-ci-output/docker-images/kubo-ci
  sed -E -i -e "/ENV GOLANG_VERSION=/s/([0-9]+\.)+[0-9]+/${go_version}/" Dockerfile

  git config --global user.email "cfcr+cibot@pivotal.io"
  git config --global user.name "CFCR CI BOT"

  set +e
  git add .
  git commit -m "Bump CI image golang to version $go_version"
  set -e
popd

