#!/bin/bash

export PS4="\t ${FUNCNAME[0]:+${FUNCNAME[0]}(): }"
set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

cp "gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

bosh_ca_cert=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/default_ca/ca)
client_secret=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/bosh_admin_client_secret)

director_ip=$(bosh-cli int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip")

"git-kubo-deployment/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" ci-service

trap "kubectl delete -f 'git-kubo-ci/specs/guestbook.yml'" EXIT

kubectl apply -f "git-kubo-ci/specs/guestbook.yml"
# wait for deployment to finish
kubectl rollout status deployment/frontend -w
kubectl rollout status deployment/redis-master -w
kubectl rollout status deployment/redis-slave -w
nodeport=$(kubectl describe svc/frontend | grep 'NodePort:' | awk '{print $3}' | sed -e 's/\/TCP//g')


worker_ip=$(BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET=${client_secret} BOSH_CA_CERT="${bosh_ca_cert}" bosh-cli -e "${director_ip}" vms | grep worker | grep -v haproxy | head -n1 | awk '{print $4}')
testvalue="$(date +%s)"

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

wget -O - "http://${worker_ip}:${nodeport}/guestbook.php?cmd=set&key=messages&value=${testvalue}"

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
