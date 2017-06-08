#!/bin/bash -ex

. "$(dirname "$0")/lib/environment.sh"

export BOSH_LOG_LEVEL=debug
export BOSH_LOG_PATH="$PWD/bosh.log"
export DEBUG=1

bosh_ca_cert=$(bosh-cli int "gcs-bosh-creds/creds.yml" --path=/default_ca/ca)
client_secret=$(bosh-cli int "gcs-bosh-creds/creds.yml" --path=/bosh_admin_client_secret)
director_ip=$(bosh-cli int "kubo-lock/metadata" --path="/internal_ip")
export WORKER_IP=$(BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET=${client_secret} BOSH_CA_CERT="${bosh_ca_cert}" bosh-cli -e ${director_ip} vms | grep worker | head -n1 | awk '{print $4}')


### vvv port to golang vvv ###
testvalue="hellothere$(date +'%N')"

if timeout 120 /bin/bash <<EOF
  until wget -O - 'http://${worker_ip}:30303/guestbook.php?cmd=set&key=messages&value=${testvalue}' | grep '{"message": "Updated"}'; do
    sleep 2
  done
EOF
then
  echo "Posted the test value to guestbook"
else
  echo "Unable to post test value to guestbook"
  exit 1
fi

wget -O - "http://${worker_ip}:30303/guestbook.php?cmd=set&key=messages&value=${testvalue}"

if timeout 120 /bin/bash <<EOF
  until wget -O - "http://${worker_ip}:30303/guestbook.php?cmd=get&key=messages" | grep ${testvalue}; do
    sleep 2
  done
EOF
then
  echo "Successfully read the test value from guestbook"
else
  echo "Expected the sample guest book to display the test value"
  exit 1
fi

kubectl delete -f "git-kubo-ci/specs/pv-guestbook.yml"
kubectl apply -f "git-kubo-ci/specs/pv-guestbook.yml"

# wait for deployment to finish
kubectl rollout status deployment/frontend -w
kubectl rollout status deployment/redis-master -w
kubectl rollout status deployment/redis-slave -w

wget -O - "http://${worker_ip}:30303/guestbook.php?cmd=set&key=messages&value=${testvalue}"

if timeout 120 /bin/bash <<EOF
  until wget -O - "http://${worker_ip}:30303/guestbook.php?cmd=get&key=messages" | grep ${testvalue}; do
    sleep 2
  done
EOF
then
  echo "Successfully read the test value from guestbook"
else
  echo "Expected the sample guest book to display the test value"
  exit 1
fi
### ^^^ port to golang ^^^ ###