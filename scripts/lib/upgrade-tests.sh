. "$DIR/lib/lb-info.sh"

query_loop() {
  local timeout_seconds=10
  local pid_to_wait="$1"
  local url="$2"
  local min_success_rate="${3:-1}"

  # output values
  local query_loop_count=0
  local query_success_count=0

  # error colors
  local red='\033[0;31m'
  local no_color='\033[0m'

  echo "Querying $url while waiting for process with pid $pid_to_wait to finish..."

  # loop while process is not finished
  while kill -0 "$pid_to_wait" >/dev/null 2>&1; do
    sleep 1

    (( query_loop_count+=1 ))

    local timestamp=`date`

    set +e

    curl -L --max-time ${timeout_seconds} -IfsS ${url} &> query_loop_last_output.txt

    if [ "$?" -ne 0 ]; then
      echo -e "[$timestamp][$query_success_count/$query_loop_count] ${red}Error: request to $url failed${no_color}"
      cat query_loop_last_output.txt
    else
      (( query_success_count+=1 ))
      if [ -z ${RUN_QUERIES_SILENTLY+x} ] || [ "$RUN_QUERIES_SILENTLY" != "1" ]; then
        echo "[$timestamp][$query_success_count/$query_loop_count] Service $url successfully responded"
      fi
    fi

    set -e
  done

  echo "$query_success_count out of $query_loop_count queries succeeded"

  local success_rate=$(echo "$query_success_count $query_loop_count" | awk '{print ($1 / $2)}')

  if (( $(echo "$success_rate $min_success_rate" | awk '{print ($1 >= $2)}') )); then
    echo "Success rate ($success_rate) passes minimum ($min_success_rate)"
  else
    echo "Success rate ($success_rate) fails to pass minimum ($min_success_rate)"
    return 1
  fi
}

wait_for_success() {
  local pid_to_wait="$1"
  local work_description="$2"

  echo "PID to wait on: $pid_to_wait, for work: $work_description"
  wait "$pid_to_wait"
  if [ "$?" -ne 0 ]; then
    echo "$work_description failed"
    exit 1
  fi
  echo "$work_description succeeded"
}

run_upgrade_test() {
  local service_name="nginx"
  local update_function="$1"
  local min_success_rate="${2:-1}"
  local component_name="${3:-"not provided"}"

  routing_mode="$(bosh-cli int environment/director.yml --path=/routing_mode)"

  if [[ "$routing_mode" == "iaas" ]]; then
    lb_address_blocking "$service_name" "$KUBO_ENVIRONMENT_DIR" "$KUBO_DEPLOYMENT_DIR"
  elif [[ "$routing_mode" == "cf" ]]; then
    cp "$PWD/kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
    cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
    "$KUBO_DEPLOYMENT_DIR/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service
    generated_service_name="$(kubectl describe service "$service_name" | grep http-route-sync | cut -d= -f2)"
    cf_apps_domain="$(bosh-cli int environment/director.yml --path=/routing_cf_app_domain_name)"
    lb_address="$generated_service_name"."$cf_apps_domain"
  else
    echo "Routing mode '$routing_mode' is not supported in this test"
    exit 1
  fi

  local lb_url="http://$lb_address"
  echo "The load balancer's URL is $lb_url"

  # update BOSH in the background
  $update_function &
  local update_pid="$!"
  echo "Update function PID: ${update_pid}"

  # exercise the load balancer URL while BOSH is updating
  local query_url="$lb_url"
  query_loop "$update_pid" "$query_url" "$min_success_rate"

  if [ "$?" != "0" ]; then
    return 1
  fi
}

upload_new_releases() {
  if [ -z ${LOCAL_DEV+x} ] || [ "$LOCAL_DEV" != "1" ]; then
    cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
    cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
  fi
  BOSH_ENV="$KUBO_ENVIRONMENT_DIR" source "$KUBO_DEPLOYMENT_DIR/bin/set_bosh_environment"

  bosh-cli upload-release https://bosh.io/d/github.com/cf-platform-eng/docker-boshrelease?v=28.0.1 --sha1 448eaa2f478dc8794933781b478fae02aa44ed6b
  bosh-cli upload-release https://github.com/pivotal-cf-experimental/kubo-etcd/releases/download/v2/kubo-etcd.2.tgz --sha1 ae95e661cd9df3bdc59ee38bf94dd98e2f280d4f
}

