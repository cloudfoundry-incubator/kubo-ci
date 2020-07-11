#!/bin/bash

set -eu -o pipefail

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

sanitize_cgroups() {
  mkdir -p /sys/fs/cgroup
  mountpoint -q /sys/fs/cgroup || \
    mount -t tmpfs -o uid=0,gid=0,mode=0755 cgroup /sys/fs/cgroup

  mount -o remount,rw /sys/fs/cgroup

  sed -e 1d /proc/cgroups | while read sys hierarchy num enabled; do
    if [ "$enabled" != "1" ]; then
      # subsystem disabled; skip
      continue
    fi

    grouping="$(cat /proc/self/cgroup | cut -d: -f2 | grep "\\<$sys\\>")"
    if [ -z "$grouping" ]; then
      # subsystem not mounted anywhere; mount it on its own
      grouping="$sys"
    fi

    mountpoint="/sys/fs/cgroup/$grouping"

    mkdir -p "$mountpoint"

    # clear out existing mount to make sure new one is read-write
    if mountpoint -q "$mountpoint"; then
      umount "$mountpoint"
    fi

    mount -n -t cgroup -o "$grouping" cgroup "$mountpoint"

    if [ "$grouping" != "$sys" ]; then
      if [ -L "/sys/fs/cgroup/$sys" ]; then
        rm "/sys/fs/cgroup/$sys"
      fi

      ln -s "$mountpoint" "/sys/fs/cgroup/$sys"
    fi
  done
}

start_docker() {
  mkdir -p /var/log
  mkdir -p /var/run

  sanitize_cgroups

  # check for /proc/sys being mounted readonly, as systemd does
  if grep '/proc/sys\s\+\w\+\s\+ro,' /proc/mounts >/dev/null; then
    mount -o remount,rw /proc/sys
  fi

  local mtu=$(cat /sys/class/net/$(ip route get 8.8.8.8|awk '{ print $5 }')/mtu)

  dockerd --data-root /scratch/docker --mtu "${mtu}" 2>/tmp/dockerd.log &
  echo $! > /tmp/docker.pid

  trap stop_docker EXIT

  sleep 1

  until docker info >/dev/null 2>&1; do
    echo waiting for docker to come up...
    sleep 1
  done
}

stop_docker() {
  local pid=$(cat /tmp/docker.pid)
  if [ -z "$pid" ]; then
    return 0
  fi

  kill -TERM $pid
}

bump_spec() {
  set -x
  tag=$(cat "$PWD/$SPEC_RELEASE_DIR/tag")

  if [[ $SPEC_BLOB_NAME == "coredns_coredns" ]]; then
    tag=$(echo $tag | sed 's/v//')
  fi

  pushd git-kubo-release-output

  local old_blob
  old_blob="$( bosh blobs --column path | grep "${SPEC_BLOB_NAME}" )"
  bosh remove-blob "$old_blob"
  scripts/download_container_images "$SPEC_IMAGE_URL:$tag"
  sed -E -i -e "/${SPEC_IMAGE_NAME}:/s/v([0-9]+\.)+[0-9]+/${tag}/" scripts/download_container_images
  find ./jobs/apply-specs/templates/specs/ -type f -exec sed -E -i -e "/${SPEC_IMAGE_NAME}:/s/v?([0-9]+\.)+[0-9]+/${tag}/" {} \;

  set +x
  cat <<EOF > "config/private.yml"
blobstore:
  options:
    json_key: '${GCS_JSON_KEY}'
EOF
  set -x

  bosh upload-blobs

  git config --global user.email "cfcr+cibot@pivotal.io"
  git config --global user.name "CFCR CI BOT"
  git add .
  git commit -m "Bump ${SPEC_NAME} to version ${tag}"
  popd
}

main() {
  start_docker
  base_dir=$PWD
  cp -r git-kubo-release/. git-kubo-release-output
  pushd git-kubo-release

  message="No new versions found to update."
  changed=""
  for i in "${!blob_names[@]}"; do
      set -x
      existing_spec_version=$(bosh blobs --column path | grep "${blob_names[i]}" | grep -o -E 'v?([0-9]+\.)+[0-9]+')
      fetched_spec_version=$(cat "$base_dir/${releases[i]}/tag")
      set +x
      if [[ "$existing_spec_version" != "$fetched_spec_version" ]] \
          && [[ "v$existing_spec_version" != "$fetched_spec_version" ]]; then
          export SPEC_RELEASE_DIR=${releases[i]}
          export SPEC_BLOB_NAME=${blob_names[i]}
          export SPEC_IMAGE_NAME=${names[i]}
          export SPEC_IMAGE_URL=${urls[i]}
          export SPEC_NAME=$(echo "${releases[i]}" | sed 's/.\{8\}$//')
          pushd "${base_dir}"
            bump_spec
          popd
          changed="$changed ${SPEC_IMAGE_NAME}"
          message="Updated releases: $changed"
      fi
  done

  echo "$message"

  popd
}
main
