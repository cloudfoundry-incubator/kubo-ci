#!/usr/bin/env bash

set -eu

export GOPATH="$PWD"
export PATH="$PATH:$GOPATH/bin"
export CONFIG_PATH="$PWD/k-drats-config/$CONFIG_PATH"

pushd src/github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests
  scripts/_run_acceptance_tests.sh
popd
