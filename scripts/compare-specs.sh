#!/bin/bash

set -exu -o pipefail

base_dir=$PWD
pushd git-kubo-release

specs=("k8s.gcr.io_heapster-amd64"
       #"k8s.gcr.io_heapster-influxdb-amd64"
       "gcr.io_google_containers_kubernetes-dashboard-amd64")
releases=("heapster-release"
          #"influxdb-release"
          "kubernetes-dashboard-release")
urls=("k8s.gcr.io/heapster-amd64"
      #"k8s.gcr.io/heapster-influxdb-amd64"
      "gcr.io/google_containers/kubernetes-dashboard-amd64")

for i in ${!specs[@]}; do
    existing_spec_version=$(bosh blobs --column path | grep "${specs[i]}" | grep -o -E 'v([0-9]+\.)+[0-9]+')
    fetched_spec_version=$(cat "$base_dir/${releases[i]}/tag")
    if [[ $existing_spec_version != $fetched_spec_version ]]; then
        cat <<EOF > "$base_dir/spec-to-update/spec.env"
export SPEC_RELEASE_DIR=${releases[i]}
export SPEC_IMAGE_NAME=${specs[i]}
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
