#!/bin/bash

set -eu

root="${PWD}"

TEST_etcd_release_version="99999+dev.$(date +%s)"
TEST_stemcell_version="$(cat "${root}/stemcell/version")"
TEST_latest_etcd_release_version="$(cat "${root}/latest-etcd-release/version")"
export TEST_etcd_release_version TEST_stemcell_version TEST_latest_etcd_release_version

consul_release_version="$(cat "${root}/consul-release/version")"
turbulence_release_version="$(cat "${root}/turbulence-release/version")"


function main() {
  setup_bosh_env_vars

  set_cloud_config
  upload_stemcell
  upload_releases
  force_compilation

  crt=$(mktemp)
  printenv BOSH_CA_CERT > "$crt"
  bosh -d eats -n deploy "${root}/etcd-release/manifests/eats.yml" \
    --vars-env=TEST \
    --var="bosh_client=${BOSH_CLIENT}" \
    --var="bosh_client_secret=${BOSH_CLIENT_SECRET}" \
    --var-file="bosh_director_ca_cert=${crt}" \
    --var="bosh_environment=https://${BOSH_ENVIRONMENT}:25555"

  bosh -d eats run-errand acceptance-tests

  bosh -n clean-up --all
}

function set_cloud_config() {
  bosh -n update-cloud-config git-kubo-ci/etcd/cloud-config.yml \
    --vars-file=kubo-lock/metadata
}

function setup_bosh_env_vars() {
  BOSH_CLIENT=bosh_admin
  BOSH_CLIENT_SECRET=$(bosh int gcs-bosh-creds/creds.yml --path=/bosh_admin_client_secret)
  BOSH_CA_CERT=$(bosh int gcs-bosh-creds/creds.yml --path=/default_ca/ca)
  BOSH_ENVIRONMENT=$(bosh int kubo-lock/metadata --path=/internal_ip)
  export BOSH_ENVIRONMENT BOSH_CA_CERT BOSH_CLIENT BOSH_CLIENT_SECRET
}

function upload_stemcell() {
  bosh upload-stemcell stemcell/stemcell.tgz
}

function upload_releases() {
  bosh upload-release turbulence-release/release.tgz
  bosh upload-release consul-release/release.tgz
  bosh upload-release latest-etcd-release/release.tgz

  bosh -n create-release --force --version "${TEST_etcd_release_version}" --dir=etcd-release
  bosh upload-release --dir=etcd-release
}

function force_compilation() {
    sed -e "s/CONSUL_RELEASE_VERSION/${consul_release_version}/g" \
      -e "s/ETCD_RELEASE_VERSION/${TEST_etcd_release_version}/g" \
      -e "s/TURBULENCE_RELEASE_VERSION/${turbulence_release_version}/g" \
      -e "s/STEMCELL_VERSION/${TEST_stemcell_version}/g" \
      "${root}/etcd-release/scripts/fixtures/eats_compilation.yml" \
      > "${root}/eats_compilation.yml"

    bosh -d compilation -n deploy "${root}/eats_compilation.yml"
    bosh -d compilation export-release "kubo-etcd/${TEST_etcd_release_version}" "ubuntu-trusty/${TEST_stemcell_version}"
    bosh -d compilation export-release "consul/${consul_release_version}" "ubuntu-trusty/${TEST_stemcell_version}"
    bosh -d compilation export-release "turbulence/${turbulence_release_version}" "ubuntu-trusty/${TEST_stemcell_version}"
    bosh -d compilation -n delete-deployment
}

function teardown() {
  set +e
  bosh -d eats -n delete-deployment
  bosh -n delete-release kubo-etcd
  bosh -n delete-release consul
  bosh -n delete-release turbulence
  set -e
}

trap teardown EXIT

main
