# query the API for the load balancer address and sets $lb_address
lb_address_blocking() {
  service_name="$1"
  kubo_environment_dir="$2"
  kubo_deployment_dir="$3"

  iaas=$(bosh-cli int "${kubo_environment_dir}/director.yml" --path="/iaas")

	lb_address=""
  current_attempt=0
  max_attempts=30

  set -e
    "${kubo_deployment_dir}/bin/set_kubeconfig" "${kubo_environment_dir}" ci-service
  set +e

  while [ -z "$lb_address" ]; do
    current_attempt=$((current_attempt+1))
    if [ ${current_attempt} -gt ${max_attempts} ]; then
      echo "Error: reached max attempts trying to obtain load balancer IP"
      return 1
    fi


    if [ ${iaas} = "gcp" ]; then
			lb_address=$(kubectl get service ${service_name} -o jsonpath={.status.loadBalancer.ingress[0].ip})
    else
      lb_address=$(kubectl get service ${service_name} -o jsonpath={.status.loadBalancer.ingress[0].hostname})
    fi

    if [ -z "$lb_address" ]; then sleep 10; fi
  done
}
