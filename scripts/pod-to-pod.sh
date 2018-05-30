#!/bin/bash

export PS4="\t ${FUNCNAME[0]:+${FUNCNAME[0]}(): }"
set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1
export DEPLOYMENT_NAME=${DEPLOYMENT_NAME:="ci-service"}

cp "gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"
cp "kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"

bosh_ca_cert=$(bosh int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/default_ca/ca)
client_secret=$(bosh int "${KUBO_ENVIRONMENT_DIR}/creds.yml" --path=/bosh_admin_client_secret)

director_ip=$(bosh int "${KUBO_ENVIRONMENT_DIR}/director.yml" --path="/internal_ip")

if [[ -f "${ROOT}/git-kubo-deployment/bin/credhub_login" ]]; then
  "${ROOT}/git-kubo-deployment/bin/credhub_login" "${KUBO_ENVIRONMENT_DIR}"
  source "${ROOT}/git-kubo-ci/scripts/lib/utils.sh"
  eval "$(set_variables)"
  "${KUBO_DEPLOYMENT_DIR}/bin/set_kubeconfig" "${cluster_name}" "${api_url}"
else
  "${KUBO_DEPLOYMENT_DIR}/bin/set_kubeconfig" "${KUBO_ENVIRONMENT_DIR}" "${DEPLOYMENT_NAME}"
fi

kubectl apply -f "git-kubo-ci/specs/pod2pod-ns.yml"

trap "kubectl -n pod2pod delete -f 'git-kubo-ci/specs/guestbook.yml'" EXIT

kubectl -n pod2pod apply -f "git-kubo-ci/specs/guestbook.yml"
# wait for deployment to finish
kubectl -n pod2pod rollout status deployment/frontend -w
kubectl -n pod2pod rollout status deployment/redis-master -w
kubectl -n pod2pod rollout status deployment/redis-slave -w
nodeport=$(kubectl -n pod2pod describe svc/frontend | grep 'NodePort:' | awk '{print $3}' | sed -e 's/\/TCP//g')


worker_ip=$(BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET=${client_secret} BOSH_CA_CERT="${bosh_ca_cert}" bosh -e "${director_ip}" vms -d "${DEPLOYMENT_NAME}"  | grep worker | head -n1 | awk '{print $4}')
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
