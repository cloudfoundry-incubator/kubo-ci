#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

metadata_path="kubo-lock/metadata"
director_name=$(bosh-cli int ${metadata_path} --path=/director_name)

mkdir -p "updated-kubo-lock"
updated_metadata_path="updated-kubo-lock/metadata"

echo "Getting creds"

credhub login
set +x

credhub get -n "/concourse/main/${director_name}/creds" --output-json | jq -r .value > bosh-creds/creds.yml

cp ${metadata_path} ${updated_metadata_path}
credhub get -n "/concourse/main/cfcr" --output-json | jq -r .value  >> ${updated_metadata_path}
credhub get -n "/concourse/main/cfcr-gcp" --output-json | jq -r .value  >> ${updated_metadata_path}
credhub get -n "/concourse/main/cfcr-gcp-cf" --output-json | jq -r .value  >> ${updated_metadata_path}
credhub get -n "/concourse/main/cfcr-gcp-cf-${director_name}-conformance" --output-json | jq -r .value >> ${updated_metadata_path}
echo "routing_mode: cf" >> ${updated_metadata_path}

set -x

credhub logout
