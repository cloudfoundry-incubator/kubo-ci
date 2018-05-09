#!/bin/bash

set -euo pipefail

creds_file=gcs-bosh-creds/creds.yml

BOSH_CA_CERT=$(bosh int ${creds_file} --path=/default_ca/ca)
BOSH_CLIENT="admin"
BOSH_CLIENT_SECRET=$(bosh int ${creds_file} --path=/admin_password)
BOSH_ENVIRONMENT=$(bosh int kubo-lock/metadata --path=/internal_ip)

export BOSH_CA_CERT BOSH_CLIENT BOSH_CLIENT_SECRET BOSH_ENVIRONMENT

