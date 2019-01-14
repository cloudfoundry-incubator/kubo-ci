#!/bin/bash
set -eux

first=$(find ./gcs-*-shipables/ -name '*shipable' | head -1)

for shipable in ./gcs-*-shipables/*shipable; do
  result=/tmp/temp_$shipable
  comm -12 "$first" "$shipable" > "$result"
  first=$result
done

new_signal_version=$(tail -1 "$result")
if [ -n "$new_signal_version" ] ; then
  echo "Found shipable version $new_signal_version" > "${SLACK_MESSAGE_FILE}"
  echo "${new_signal_version}" > "${SHIPABLE_VERSION_FILE}"
  exit 0
fi
echo "Failed to find a shipable tarball" > "${SLACK_MESSAGE_FILE}"
