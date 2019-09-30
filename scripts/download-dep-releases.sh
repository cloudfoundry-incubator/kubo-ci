#!/bin/bash

set -euo pipefail

main() {
  cd dep-releases
  yq read ../git-kubo-deployment/manifests/cfcr.yml releases.*.url | grep -v null | grep -v kubo-[0-9] | sed 's|^-\ ||g' | xargs -n 1 curl -SLJO
  ls -alh
}

