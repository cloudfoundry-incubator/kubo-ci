#!/bin/bash
set -eux

first=$(find ./gcs-*-shipables/ -name '*shipable' | head -1)

for shipable in ./gcs-*-shipables/*shipable; do
  result=/tmp/temp_$(basename "$shipable")
  comm -12 "$first" "$shipable" > "$result"
  first=$result
done

new_signal_version=$(tail -1 "$result")
if [ -n "$new_signal_version" ] ; then
  release_sha=$(echo "$new_signal_version" | awk -F' ' '{print $1}')
  deployment_sha=$(echo "$new_signal_version" | awk -F' ' '{print $2}')
  echo "Ready to :ship: <https://github.com/cloudfoundry-incubator/kubo-release/tree/${release_sha}|${release_sha}> <https://github.com/cloudfoundry-incubator/kubo-deployment/tree/${deployment_sha}/|${deployment_sha}>" > "${SLACK_MESSAGE_FILE}"
  echo "${new_signal_version}" > "${SHIPABLE_VERSION_FILE}"
  exit 0
fi
echo "Failed to find a shipable tarball" > "${SLACK_MESSAGE_FILE}"
exit 1
