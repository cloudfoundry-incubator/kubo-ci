#!/bin/sh

export KUBO_DEPLOYMENT_DIR="${KUBO_DEPLOYMENT_LOCATION:-"git-kubo-deployment"}"
export KUBO_ENVIRONMENT_DIR="${PWD}/environment"
mkdir -p "${KUBO_ENVIRONMENT_DIR}"
