#!/bin/bash

set -eo pipefail
set -x

sudo add-apt-repository ppa:gophers/archive
sudo apt update
sudo apt-get install -y golang-1.9-go

go get github.com/aquasecurity/kube-bench
cp $GOROOT/bin/kube-bench .
./kube-bench help
