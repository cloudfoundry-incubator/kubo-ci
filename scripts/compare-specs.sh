#!/bin/bash

set -exu -o pipefail

base_dir=$PWD
pushd git-kubo-release

blob_names=("k8s.gcr.io_kubernetes-dashboard-amd64"
            "k8s.gcr.io_metrics-server-amd64"
            "coredns_coredns")
releases=("kubernetes-dashboard-release"
          "metrics-server-release"
          "coredns-release")
urls=("k8s.gcr.io/kubernetes-dashboard-amd64"
      "k8s.gcr.io/metrics-server-amd64"
      "coredns/coredns")
names=("kubernetes-dashboard-amd64"
       "metrics-server-amd64"
       "coredns")

for i in ${!blob_names[@]}; do
    existing_spec_version=$(bosh blobs --column path | grep "${blob_names[i]}" | grep -o -E 'v?([0-9]+\.)+[0-9]+')
    fetched_spec_version=$(cat "$base_dir/${releases[i]}/tag")
    if [[ $existing_spec_version != $fetched_spec_version ]]; then
        cat <<EOF > "$base_dir/spec-to-update/spec.env"
export SPEC_RELEASE_DIR=${releases[i]}
export SPEC_BLOB_NAME=${blob_names[i]}
export SPEC_IMAGE_NAME=${names[i]}
export SPEC_IMAGE_URL=${urls[i]}
export SPEC_NAME=$(echo "${releases[i]}" | sed 's/.\{8\}$//')
EOF
        break
    fi
done

popd

if [ ! -f spec-to-update/spec.env ]; then
    echo "No new versions found to update."
    exit 0
fi
