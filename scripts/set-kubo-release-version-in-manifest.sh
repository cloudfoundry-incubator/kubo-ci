#! /bin/bash

cp -r git-kubo-deployment/. git-kubo-deployment-output

if [[ $IS_FINAL ]]; then
cat << EOF > replace-kubo-version.yml
- type: replace
  path: /releases/name=kubo
  value:
    name: kubo
    version: ((version))
    sha1: ((sha))
    url: ((url))
EOF
version="$(cat kubo-version/version)"
sha="$(shasum kubo-release-tarball/kubo-release-${version}.tgz | cut -d ' ' -f 1)"
url="https://github.com/cloudfoundry-incubator/kubo-release/releases/download/v${version}/kubo-release-${version}.tgz"
bosh int git-kubo-deployment/manifests/cfcr.yml -o replace-kubo-version.yml -v version="$version" -v sha="$sha" -v url="$url" > git-kubo-deployment-output/manifests/cfcr.yml
else
cat << EOF > replace-kubo-version.yml
- type: replace
  path: /releases/name=kubo
  value:
    name: kubo
    version: latest
EOF
bosh int git-kubo-deployment/manifests/cfcr.yml -o replace-kubo-version.yml > git-kubo-deployment-output/manifests/cfcr.yml
fi
