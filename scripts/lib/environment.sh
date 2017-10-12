#!/bin/sh

export KUBO_DEPLOYMENT_DIR="${KUBO_DEPLOYMENT_LOCATION:-"git-kubo-deployment"}"

export KUBO_CI_DIR="${KUBO_CI_DIR:-"git-kubo-ci"}"
if [ ! -d "$KUBO_CI_DIR" ]; then
  echo "KUBO_CI_DIR $KUBO_CI_DIR does not exists"
  exit 1
fi

if ([ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]) || [ -z "$KUBO_ENVIRONMENT_DIR" ]; then
  export KUBO_ENVIRONMENT_DIR="${PWD}/environment"
  mkdir -p "${KUBO_ENVIRONMENT_DIR}"
fi
