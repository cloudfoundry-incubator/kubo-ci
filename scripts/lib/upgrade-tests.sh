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

    curl -L --max-time ${timeout_seconds} -IfsS ${url} &> query_loop_last_output.txt

    if [ "$?" -ne 0 ]; then
      echo -e "[$timestamp][$query_success_count/$query_loop_count] ${red}Error: request to $url failed${no_color}"
      cat query_loop_last_output.txt
    else
      (( query_success_count+=1 ))
      echo "[$timestamp][$query_success_count/$query_loop_count] Service $url successfully responded"
    fi
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

  echo "PID to wait on: $pid_to_wait"
  wait "$pid_to_wait"
  if [ "$?" -ne 0 ]; then
    echo "$work_description failed"
    return 1
  fi

  echo "$work_description succeeded"
}

run_upgrade_test() {
  local service_name="nginx"
  local update_function="$1"
  local min_success_rate="${2:-1}"

  lb_address_blocking "$service_name" "$KUBO_ENVIRONMENT_DIR" "$KUBO_DEPLOYMENT_DIR"
  if [ -z "$lb_address" ]; then
    echo "Error: couldn't obtain load balancer address"
    return 1
  fi

  local lb_url="http://$lb_address"
  echo "The load balancer's URL is $lb_url"

  # update BOSH in the background
  $update_function &
  local update_pid="$!"

  # exercise the load balancer URL while BOSH is updating
  local query_url="$lb_url"
  query_loop "$update_pid" "$query_url" "$min_success_rate" &
  local query_loop_pid="$!"

  wait_for_success "$update_pid" "Update"
  update_code="$?"

  wait_for_success "$query_loop_pid" "HA query loop"
  query_code="$?"

  if [ "$update_code" != "0" ] || [ "$query_code" != "0" ]; then
    return 1
  fi
}
