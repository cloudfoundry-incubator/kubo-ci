#!/bin/bash

set -exu -o pipefail

check_for_raised_pr() {
  output=$(curl -u cfcr:${CFCR_USER_TOKEN} -H "Content-Type: application/json" -X GET https://api.github.com/repos/cloudfoundry-incubator/kubo-release/issues\?creator\=CFCR | jq '.[] | select (. | has("pull_request")) | select (.title | contains("Flannel upgrade"))' | wc -c)
  if [ $output -ne 0 ]; then return 0; else return 1; fi
}



fetched_version=$(cat "$PWD/flannel-release/version")
existing_version=$(grep -o -E "([0-9]+\.)+[0-9]+" git-kubo-release/packages/flanneld/packaging)

if [ "$fetched_version" == "$existing_version" ]; then
  echo "No new flannel versions fetched"
  exit 1
fi

if check_for_raised_pr ; then
  echo "Upgrade PR for this version has already been raised"
  exit 1
fi

