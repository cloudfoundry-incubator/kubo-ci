#! /usr/bin/env bash
set -ex

credhub_login() {
  local environment="$1"

  local credhub_admin_secret=$(bosh int "${environment}/creds.yml" --path="/credhub_admin_client_secret")
  local credhub_api_url="https://$(bosh int "${environment}/director.yml" --path="/internal_ip"):8844"

  credhub login --client-name credhub-admin --client-secret "${credhub_admin_secret}" \
    -s "${credhub_api_url}" \
    --ca-cert <(bosh int "${environment}/creds.yml" --path="/credhub_tls/ca") \
    --ca-cert <(bosh int "${environment}/creds.yml" --path="/uaa_ssl/ca") 1>/dev/null
}

main() {
  local environment=$1
  local deployment=$2

  credhub_login ${environment}

  local project_id=$(bosh int ${environment}/director.yml --path=/project_id)
  local prefix=$(bosh int ${environment}/director.yml --path=/director_name)
  local region=$(bosh int ${environment}/director.yml --path=/zone | rev | cut -c 3- | rev)
  local service_account_key=$(bosh int ${environment}/director.yml --path=/gcp_service_account)
  local credhub_output=$(credhub get -n "/${environment}/${deployment}/tls-kubernetes")
  local ca=$(bosh int <(echo $credhub_output) --path=/value/ca)
  local certificate=$(bosh int <(echo $credhub_output) --path=/value/certificate)
  local private_key=$(bosh int <(echo $credhub_output) --path=/value/private_key)

  local cert_chain=$(echo -e "$certificate\n$ca")

  terraform plan \
    -var prefix=${deployment:-cfcr} \
    -var project_id=${project_id} \
    -var region=${region} \
    -var service_account_key_path=<(echo ${service_account_key}) \
    -var private_key_path=<(echo ${private_key}) \
    -var certificate_key_path=<(echo ${cert_chain}) \
    gcp-https-lb.tf

#  terraform apply \
#    -var prefix=${deployment:-cfcr} \
#    -var project_id=${project_id} \
#    -var region=${region} \
#    -var service_account_key_path=<(echo ${service_account_key}) \
#    -var private_key_path=<(echo ${private_key}) \
#    -var certificate_key_path=<(echo ${cert_chain}) \
#    gcp-https-lb.tf
}

main $@
