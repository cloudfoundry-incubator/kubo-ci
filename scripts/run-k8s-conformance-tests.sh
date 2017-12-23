#!/usr/bin/env bash

[ -z "$DEBUG" ] || set -x

set -eu
set -o pipefail

BASE_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)

verify_args() {
  set +e # Cant be set since read returns a non-zero when it reaches EOF
  read -r -d '' usage <<-EOF
	Usage: $(basename "$0") [-h] environment deployment-name release-tarball results-dir
	
	Help Options:
		-h  show this help text
	EOF
  set -e

  while getopts ':h:' option; do
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
  if [[ $# -lt 3 ]]; then
    echo "$usage" >&2
    exit 64
  fi
}

run_tests() {
  local environment="$1"
  local deployment="$2"
  local results_dir="$3"

  local iaas=$(bosh-cli int "$environment/director.yml" --path='/iaas')
  local routing_mode=$(bosh-cli int "$environment/director.yml" --path='/routing_mode')
  local director_name=$(bosh-cli int "$environment/director.yml" --path='/director_name')

  if [ ! -d "$results_dir" ]; then
    echo "Error: $results_dir does not exist" >&2
    exit 1
  fi

  local bosh_name=$(basename "$environment")
  local release_version=$(bosh-cli int <(bosh-cli -e "$bosh_name" -d "$deployment" manifest) --path='/releases/name=kubo/version')
  if [ -z "$release_version" ]; then
    echo "Error: Unable to retrieve release version from BOSH" >&2
    exit 1
  fi

  local tmpfile=$(mktemp)
  $BASE_DIR/scripts/generate-test-config.sh $environment $deployment > $tmpfile
  export CONFIG=$tmpfile

  export CONFORMANCE_RESULTS_DIR="$results_dir"
  export CONFORMANCE_RELEASE_VERSION="$release_version"
  ginkgo -progress -v "$BASE_DIR/src/tests/conformance"

  return 0
}

main() {
  verify_args "$@"
  run_tests "$@"
}

main "$@"
