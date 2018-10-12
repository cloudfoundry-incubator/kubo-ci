#!/usr/bin/env bash

set -eu -o pipefail

print_usage() {
    echo "usage: $(basename ${0}) K8S_VERSION K8S_PATH SPEC_PATH"
    echo "example: $(basename ${0}) v1.11.3 ~/go/src/k8s.io/kubernetes ~/workspace/kubo-release/jobs/kube-apiserver/spec"
}

if [[ $# -ne 3 ]]; then
    print_usage
    exit 1
fi

K8S_VERSION=$1
shift
K8S_PATH=$1
shift
SPEC_PATH=$1

pushd "$K8S_PATH"
    git fetch upstream --tags
    git checkout "$K8S_VERSION"
popd

JOBSPEC_DIRECTORY=$K8S_PATH/jobspec
trap "rm -rf $JOBSPEC_DIRECTORY" EXIT
mkdir "$JOBSPEC_DIRECTORY"
cp main.go "$JOBSPEC_DIRECTORY"
cp -r flag_generator "$JOBSPEC_DIRECTORY"

unset GOPATH

pushd "$JOBSPEC_DIRECTORY"
    go run main.go "$SPEC_PATH"
popd
