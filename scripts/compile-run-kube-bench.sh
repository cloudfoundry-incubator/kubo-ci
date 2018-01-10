#!/bin/bash

set -eo pipefail

print_usage() {
  echo "Usage: compile-run-kube-bench.sh master|node"
}

if [ -z "$1" ]; then
  print_usage
  exit 1
elif [ "$1" != "master" ] && [ "$1" != "node" ]; then
  print_usage
  exit 1
else
  node_type=$1
fi

set -x

sudo add-apt-repository -y ppa:gophers/archive
sudo apt update
sudo apt-get install -y golang-1.9-go git

GOROOT="/usr/lib/go-1.9"
export GOROOT

"$GOROOT/bin/go" get github.com/aquasecurity/kube-bench

cp -R "$HOME/go/src/github.com/aquasecurity/kube-bench/cfg" .

config_path="$PWD/kube-bench-config.yml"

# kubectl needs to be in the PATH
PATH="$PATH:/var/vcap/packages/kubernetes/bin"
export PATH

# kubectl version (used by kube-bench) needs a kubeconfig
if [ "$node_type" == "master" ]; then
  KUBECONFIG="/var/vcap/jobs/kube-controller-manager/config/kubeconfig"
else
  KUBECONFIG="/var/vcap/jobs/kubelet/config/kubeconfig"
fi
export KUBECONFIG

~/go/bin/kube-bench \
  --config="$config_path" \
  "$node_type"

