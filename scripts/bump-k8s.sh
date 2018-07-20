#!/bin/bash
set -exu -o pipefail

source git-kubo-ci/scripts/lib/generate-pr.sh

pr_kubo_ci() {
  version="$1"
  tag="$2"
  docker_file="docker-images/conformance/Dockerfile"

  cp -r git-kubo-ci/. git-kubo-ci-output
  pushd git-kubo-ci-output

  existing_k8s_version="$(grep "ENV KUBE_VERSION" "$docker_file" | sed 's/.*"v\(.*\).*"/\1/g')"

  if [ $existing_k8s_version == $version ]; then
      echo "Kubernetes version already up-to-date."
      exit 0
  else
    sed -i "s/ENV KUBE_VERSION=.*/ENV KUBE_VERSION=\"v${version}\"/g" "$docker_file"
    generate_pull_request "k8s_in_conformance" $tag
  fi

  popd
}

pr_kubo_release() {
  version="$1"
  tag="$2"
  name="kubernetes-${version}"

  cp -r git-kubo-release/. git-kubo-release-output
  pushd git-kubo-release-output

  existing_k8s_spec=$(bosh blobs --column path | grep kubelet | grep -o -E 'kubernetes-([0-9]+\.)+[0-9]+')
  existing_k8s_version=$(echo $existing_k8s_spec | grep -o -E '([0-9]+\.)+[0-9]+')

  if [ $existing_k8s_version == $version ]; then
      echo "Kubernetes version already up-to-date."
      exit 0
  else
      ./scripts/download_k8s_binaries $version
      pushd packages/kubernetes
      sed -E -i -e "s/([0-9]+\.)+[0-9]+/${version}/" packaging
      sed -E -i -e "s/${existing_k8s_spec}/${name}/" spec
      popd

      cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF
      bosh upload-blobs
      generate_pull_request "kubernetes" $tag
  fi

  popd
}

tag=$(cat "$PWD/k8s-release/tag")
version=$(cat "$PWD/k8s-release/version")

pr_kubo_release "$version" "$tag"
pr_kubo_ci "$version" "$tag"



