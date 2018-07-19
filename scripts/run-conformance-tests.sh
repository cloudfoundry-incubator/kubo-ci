#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

run_conformance_test() {
  ginkgo --flakeAttempts=2 --nodes=8 -p -progress -focus  "\[Conformance\]" -skip "\[Serial\]" /e2e.test
  ginkgo --flakeAttempts=2 -focus="\[Serial\].*\[Conformance\]" /e2e.test
}

main() {
  cp gcs-kubeconfig/config $KUBECONFIG

  run_conformance_test
}

main
