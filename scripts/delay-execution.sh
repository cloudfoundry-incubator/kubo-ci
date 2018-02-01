#!/bin/bash

set -exu -o pipefail

echo "Delaying execution by ${DELAY_TIME_SECS} seconds"
sleep "${DELAY_TIME_SECS}"
