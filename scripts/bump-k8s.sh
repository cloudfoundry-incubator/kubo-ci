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

pr_release() {
  version="$1"
  tag="$2"
  release_name="$3"

  git_release_name="git-${release_name}"

  cp -r "${git_release_name}/." "${git_release_name}-output"
  pushd "${git_release_name}-output"

  ./scripts/download_k8s_binaries $version

  if [ -n "$(git status --porcelain)" ]; then
    cat <<EOF > "config/private.yml"
blobstore:
  options:
    credentials_source: static
    json_key: '${GCS_JSON_KEY}'
EOF
    bosh upload-blobs
    generate_pull_request "kubernetes" "$tag" "${release_name}" "develop"
  else
    echo "Kubernetes version is already up-to-date"
  fi

  popd
}

tag=$(cat "$PWD/k8s-release/tag")
version=$(cat "$PWD/k8s-release/version")

if [ "${REPO:-}" == "ci" ]; then
  pr_kubo_ci "$version" "$tag"
elif [ "${REPO:-}" == "windows" ]; then
  pr_release "$version" "$tag" "kubo-release-windows"
else
  pr_release "$version" "$tag" "kubo-release"
fi
