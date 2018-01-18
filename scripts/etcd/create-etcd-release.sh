#!/bin/bash

set -exu -o pipefail

cd git-kubo-etcd

bosh create-release --sha2 --tarball="../etcd-release/etcd-release.tgz"
