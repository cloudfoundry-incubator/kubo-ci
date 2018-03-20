#!/bin/bash

SEMVER_REGEX="^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(\.(0|[1-9][0-9]*))?$"

semver_arr() {
  version=$1
  if [[ "$version" =~ $SEMVER_REGEX ]]; then
    local major=${BASH_REMATCH[1]}
    local minor=${BASH_REMATCH[2]}
    local patch=${BASH_REMATCH[3]}
    if [ ! -z $patch ]; then
      patch=$(echo $patch | cut -c 2-)
    fi
    eval "$2=(\"$major\" \"$minor\" \"$patch\")"
  fi
}

compare_semvers() {
  semver_arr $1 a
  semver_arr $2 b

  for i in 0 1 2; do
    local x=${a[$i]}
    local y=${b[$i]}
    if [ -z $x ]; then
      x=0
    fi
    if [ -z $y ]; then
      y=0
    fi
    local diff=$(($x - $y))
    if [[ $diff -lt 0 ]]; then
      echo -1; return 0
    elif [[ $diff -gt 0 ]]; then
      echo 1; return 0
    fi
  done
  echo 0
}
