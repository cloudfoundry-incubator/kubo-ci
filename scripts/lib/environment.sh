#!/bin/sh

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/../../.." && pwd )"

export KUBO_DEPLOYMENT_DIR="${KUBO_DEPLOYMENT_LOCATION:-"${ROOT}/git-kubo-deployment"}"

export KUBO_CI_DIR="${KUBO_CI_DIR:-"${ROOT}/git-kubo-ci"}"
if [ ! -d "$KUBO_CI_DIR" ]; then
  echo "KUBO_CI_DIR $KUBO_CI_DIR does not exists"
  exit 1
fi

if ([ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]) || [ -z "$KUBO_ENVIRONMENT_DIR" ]; then
  export KUBO_ENVIRONMENT_DIR="${ROOT}/environment"
  mkdir -p "${KUBO_ENVIRONMENT_DIR}"
fi
