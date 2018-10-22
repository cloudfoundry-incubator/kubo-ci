#!/usr/bin/env bash

set -eu -o pipefail

print_usage() {
    echo "usage: $(basename ${0}) K8S_VERSION K8S_PATH K8S_COMPONENT SPEC_PATH"
    echo "example: $(basename ${0}) v1.11.3 ~/go/src/k8s.io/kubernetes kube-apiserver ~/workspace/kubo-release/jobs/kube-apiserver/spec"
}

if [[ $# -ne 4 ]]; then
    print_usage
    exit 1
fi

K8S_VERSION=$1
shift
K8S_PATH=$1
shift
K8S_COMPONENT=$1
shift
SPEC_PATH=$1

pushd "$K8S_PATH"
    git fetch origin --tags
    git checkout "$K8S_VERSION"
popd

JOBSPEC_DIRECTORY=$K8S_PATH/jobspec
trap "rm -rf $JOBSPEC_DIRECTORY" EXIT
mkdir "$JOBSPEC_DIRECTORY"
cp main.go "$JOBSPEC_DIRECTORY"
mkdir "$JOBSPEC_DIRECTORY/flag_generator"
cp -r flag_generator/*.go "$JOBSPEC_DIRECTORY/flag_generator"

unset GOPATH

pushd "$JOBSPEC_DIRECTORY"
    go run main.go "$K8S_COMPONENT" "$SPEC_PATH"
popd
