#!/bin/bash
set -euox pipefail
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)

export KUBECONFIG=$ROOT/gcs-kubeconfig/config

curl -sL https://istio.io/downloadIstioctl | sh -

export PATH=$PATH:$HOME/.istioctl/bin

istioctl install --set profile=demo -y
