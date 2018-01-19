#!/bin/bash

set -exu -o pipefail

cd release-dir

bosh create-release --timestamp-version --sha2 --tarball="../release/release.tgz"
