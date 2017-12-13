#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

metadata_path="kubo-lock/metadata"
director_name=$(bosh-cli int ${metadata_path} --path=/director_name)

echo "Getting creds"

credhub login
set +x

credhub get -n "/concourse/main/${director_name}/creds" --output-json | jq .value > bosh-creds/creds.yml

set -x
credhub logout
