#!/usr/bin/env bash

set -e

ENV_TYPE="$env"

if [ -z "$ENV_TYPE" ]; then
  echo "Environment variable 'env' must be set to gcp, gcp-lb"
  exit 1
fi

GIT_KUBO_DEPLOYMENT=git-kubo-deployment

generate_director_yml_template() {
  case "$ENV_TYPE" in
    gcp)
      ops_file_args="-o "$GIT_KUBO_DEPLOYMENT"/environment-configuration/gcp/gcp.yml \
        -o "$GIT_KUBO_DEPLOYMENT"/environment-configuration/authorization/rbac.yml \
        -o "$GIT_KUBO_DEPLOYMENT"/environment-configuration/routing/cf-routing.yml"
      ;;
    gcp-lb)
      ops_file_args="-o "$GIT_KUBO_DEPLOYMENT"/environment-configuration/gcp/gcp.yml \
        -o "$GIT_KUBO_DEPLOYMENT"/environment-configuration/authorization/rbac.yml \
        -o "$GIT_KUBO_DEPLOYMENT"/environment-configuration/routing/iaas-routing.yml"
      ;;
    *)
      echo "${ENV_TYPE} is not yet supported"
      exit 1
      ;;
  esac

  bosh interpolate "$GIT_KUBO_DEPLOYMENT/environment-configuration/director.yml" $ops_file_args
}

generate_vars_file() {
  tmpfile=$(mktemp)
  cat $HOME/workspace/kubo-locks/gcp-with-bosh/unclaimed/crow > "$tmpfile"
  echo "$IAAS_CREDENTIALS" >> "$tmpfile"
  #global_variables $env >> "$tmpfile"
  cat "$tmpfile"
}

generate_director_yml() {
  director_template_path="$1"
  vars_file_path="$2"

  bosh interpolate "$director_template_path" --vars-file "$vars_file_path"
}

main() {
  generate_vars_file > vars.yml
  generate_director_yml_template > director-template.yml
  generate_director_yml director-template.yml vars.yml > generated-config/director.yml
}

main
