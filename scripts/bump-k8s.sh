#!/bin/bash
set -exu -o pipefail

source git-kubo-ci/scripts/lib/generate-pr.sh

pr_kubo_ci() {
  version="$1"
  tag="$2"
  cp -r git-kubo-ci/. git-kubo-ci-output
  pushd git-kubo-ci-output

  docker_conformance_file="docker-images/conformance/Dockerfile"
  existing_conformance_k8s_version="$(grep "ENV KUBE_VERSION" "$docker_conformance_file" | sed 's/.*"v\(.*\).*"/\1/g')"

  if [ "$existing_conformance_k8s_version" != "$version" ]; then
      sed -i "s/ENV KUBE_VERSION=.*/ENV KUBE_VERSION=\"v${version}\"/g" "$docker_conformance_file"
  fi

  docker_ci_file="docker-images/kubo-ci/Dockerfile"
  existing_ci_k8s_version="$(grep "ENV KUBE_VERSION" "$docker_ci_file" | sed 's/.*"v\(.*\).*"/\1/g')"

  if [ "$existing_ci_k8s_version" != "$version" ]; then
      sed -i "s/ENV KUBE_VERSION=.*/ENV KUBE_VERSION=\"v${version}\"/g" "$docker_ci_file"
  fi

  if [ "$existing_conformance_k8s_version" == "$version" ] && [ "$existing_ci_k8s_version" == "$version" ]; then
      echo "Kubernetes version already up-to-date."
  else
      generate_pull_request "k8s" "$tag" "kubo-ci" "master"
  fi
  popd
}

pr_kubo_release() {
  version="$1"
  tag="$2"

  cp -r git-kubo-release/. git-kubo-release-output
  pushd git-kubo-release-output

  ./scripts/download_k8s_binaries $version

  if [ -n "$(git status --porcelain)" ]; then
    cat <<EOF > "config/private.yml"
blobstore:
  options:
    access_key_id: ${ACCESS_KEY_ID}
    secret_access_key: ${SECRET_ACCESS_KEY}
EOF
    bosh upload-blobs
    generate_pull_request "kubernetes" "$tag" "kubo-release" "develop"
  else
    echo "Kubernetes version is already up-to-date"
  fi

  popd
}

tag=$(cat "$PWD/k8s-release/tag")
version=$(cat "$PWD/k8s-release/version")

pr_kubo_release "$version" "$tag"
pr_kubo_ci "$version" "$tag"
