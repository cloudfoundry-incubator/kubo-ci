#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

BASE_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)

verify_args() {
  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' usage <<-EOF
	Usage: $(basename "$0") [-h] environment deployment-name
	
	Help Options:
		-h  show this help text
	EOF
  set -e

  while getopts ':h' option; do
    case "$option" in
      h) echo "$usage"
         exit 0
         ;;
     \?) printf "Illegal option: -%s\n" "$OPTARG" >&2
         echo "$usage" >&2
         exit 64
         ;;
    esac
  done
  shift $((OPTIND - 1))
  if [[ $# -lt 2 ]]; then
    echo "$usage" >&2
    exit 64
  fi
}

run_tests() {
  local environment="$1"
  local deployment="$2"

  local iaas=$(bosh-cli int "$environment/director.yml" --path='/iaas')

  local tmpfile=$(mktemp)
  $BASE_DIR/scripts/generate-test-config.sh $environment $deployment > $tmpfile
  export CONFIG=$tmpfile

  ginkgo -progress -v "$BASE_DIR/src/tests/turbulence-tests/worker_failure"
  ginkgo -progress -v "$BASE_DIR/src/tests/turbulence-tests/master_failure"
  if [[ "${iaas}" == "gcp" || "${iaas}" == "aws" || "${iaas}" == "vsphere" ]]; then
    ginkgo -progress -v "$BASE_DIR/src/tests/turbulence-tests/persistence_failure"
  fi

  return 0
}

main() {
  verify_args "$@"
  run_tests "$@"
}

main "$@"
