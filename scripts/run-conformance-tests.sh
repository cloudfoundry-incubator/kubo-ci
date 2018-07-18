#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

run_conformance_test() {
  ginkgo -p -progress -focus  "\[Conformance\]" -skip "\[Serial\]" /e2e.test
  ginkgo -focus="\[Serial\].*\[Conformance\]" /e2e.test
}

main() {
  cp gcs-kubeconfig/config $KUBECONFIG

  run_conformance_test
}

main
