#!/bin/sh
set -euo pipefail

cleanup_service() {
  local space_guid=$1
  local service_name=$2
  local guid=$(cf curl "/v2/spaces/${space_guid}/service_instances?q=name%3A${service_name}" | jq -r '.resources[0].metadata.guid')

  echo "Deleting Service ${service_name}"

  cf curl -X DELETE "/v2/service_instances/${guid}?recursive=true&accepts_incomplete=true&purge=true"
}

cleanup_space() {
   local org_guid=$1
   local space=$2
   local guid=$(cf curl "/v2/spaces?q=name%3A${space}&q=organization_guid%3A${org_guid}" | jq -r '.resources[0].metadata.guid')
   local service_name=$(cf curl "/v2/spaces/${guid}/summary" | jq -r '.services[].name')

   for_each2 cleanup_service "$guid" $(cf curl "/v2/spaces/${guid}/summary" | jq -r '.services[].name')
}

cleanup_org() {
   local org=$1
   local guid=$(cf org "${org}" --guid)

   for_each2 cleanup_space "$guid" $(cf curl "/v2/organizations/${guid}/spaces" | jq -r '.resources[].entity.name')
   cf delete-org -f "$org"
}

for_each2() {
  local f=$1; shift
  local arg1=$1; shift
  for arg2 in $@; do
    "$f" "$arg1" "$arg2"
  done
}

for_each() {
  local f=$1; shift
  for arg in $@; do
    "$f" "$arg"
  done
}

cleanup_cf() {
  local cf_api_url=$1
  local cf_admin_password=$2
  local env_name=$3

  cf api "${cf_api_url}" --skip-ssl-validation
  cf auth admin "${cf_admin_password}"
  for_each cleanup_org $(cf curl /v2/organizations?order-by=name | jq -r ".resources[].entity.name"  | grep "${env_name}")
}

main() {
    local cf_api_url=$(bosh-cli int "$ENV_FILE" --path=/routing_cf_api_url)
    cleanup_cf "$cf_api_url" "$CF_PASSWORD" "$ENV_NAME"
}

main $@