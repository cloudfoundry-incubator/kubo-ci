#!/bin/bash

set -exu -o pipefail

GO_VERSION=$(cat $PWD/golang-version/component-golang-version)

cp -r git-kubo-ci/. git-kubo-ci-output

# pushd golang-release
#   blob_name=$(bosh blobs --json | jq '.Tables[0].Rows[] | .path | select(test("'"${GO_VERSION}"'.*linux"))' --raw-output)
#   with_prefix="${blob_name%.linux-amd64.tar.gz}"
#   go_version="${with_prefix#go}"
# popd

pushd git-kubo-ci-output/docker-images/kubo-ci
  sed -E -i -e "/ENV GOLANG_VERSION=/s/([0-9]+\.)+[0-9]+/${GO_VERSION}/" Dockerfile

  git config --global user.email "cfcr+cibot@pivotal.io"
  git config --global user.name "CFCR CI BOT"

  set +e
  git add .
  git commit -m "Bump CI image golang to version $GO_VERSION"
  set -e
popd

