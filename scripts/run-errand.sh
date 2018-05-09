#!/bin/bash

set -euo pipefail

source "$(dirname "$0")/lib/environment.sh"

bosh run-errand --instance "${INSTANCE}" "${ERRAND_NAME}" -n --keep-alive
