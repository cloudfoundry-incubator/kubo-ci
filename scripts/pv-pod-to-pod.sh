#!/bin/bash

set -exu -o pipefail

deploy_guestbook() {
  kubectl apply -f "git-kubo-ci/specs/pv-guestbook.yml"
  # wait for deployment to finish
  kubectl rollout status deployment/frontend -w
  kubectl rollout status deployment/redis-master -w
  kubectl rollout status deployment/redis-slave -w
  nodeport=$(kubectl describe svc/frontend | grep 'NodePort:' | awk '{print $3}' | sed -e 's/\/TCP//g')
}

delete_guestbook() {
  kubectl delete -f "git-kubo-ci/specs/pv-guestbook.yml" --timeout=0
}

post_to_guestbook() {
  if timeout 120 /bin/bash <<EOF
    until wget -O - 'http://${worker_ip}:${nodeport}/guestbook.php?cmd=set&key=messages&value=${testvalue}' | grep '{"message": "Updated"}'; do
      sleep 2
    done
EOF
  then
    echo "Posted the test value to guestbook"
  else
    echo "Unable to post test value to guestbook"
    exit 1
  fi
}

get_from_guestbook() {
  if timeout 120 /bin/bash <<EOF
    until wget -O - "http://${worker_ip}:${nodeport}/guestbook.php?cmd=get&key=messages" | grep ${testvalue}; do
      sleep 2
    done
EOF
  then
    echo "Successfully read the test value from guestbook"
  else
    echo "Expected the sample guest book to display the test value"
    exit 1
  fi
}

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

. "$(dirname "$0")/lib/environment.sh"

cp "gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

iaas=$(bosh-cli int "kubo-lock/metadata" --path=/iaas)
director_ip=$(bosh-cli int "kubo-lock/metadata" --path="/internal_ip")
bosh_ca_cert=$(bosh-cli int "gcs-bosh-creds/creds.yml" --path=/default_ca/ca)
client_secret=$(bosh-cli int "gcs-bosh-creds/creds.yml" --path=/bosh_admin_client_secret)
worker_ip=$(BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET=${client_secret} BOSH_CA_CERT="${bosh_ca_cert}" bosh-cli -e ${director_ip} vms | grep worker | head -n1 | awk '{print $4}')

testvalue="$(date +%s)"

"git-kubo-deployment/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service

if [ -e "git-kubo-ci/specs/storage-class-${iaas}.yml" ]; then 
  kubectl create -f "git-kubo-ci/specs/storage-class-${iaas}.yml"
  kubectl create -f "git-kubo-ci/specs/persistent-volume-claim.yml"
  deploy_guestbook
  post_to_guestbook
  get_from_guestbook
  delete_guestbook
  deploy_guestbook
  get_from_guestbook
  delete_guestbook
  kubectl delete -f "git-kubo-ci/specs/persistent-volume-claim.yml"
  kubectl delete -f "git-kubo-ci/specs/storage-class-${iaas}.yml"
else
  echo "Skipping test as no storage-class-${iaas}.yml file exists"
fi
