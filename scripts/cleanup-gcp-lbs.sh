#!/bin/bash

set -eux

export GOOGLE_APPLICATION_CREDENTIALS="$(mktemp)"

set +x
echo "$GCP_SERVICE_ACCOUNT" > "$GOOGLE_APPLICATION_CREDENTIALS"
set -x

gcloud auth activate-service-account --key-file "$GOOGLE_APPLICATION_CREDENTIALS" --project cf-pcf-kubo

forwardingRules=""
targetPools=""

for lb in $(gcloud compute target-pools list | awk '/^a/ {printf "https://www.googleapis.com/compute/v1/projects/cf-pcf-kubo/regions/%s/targetPools/%s\n", $2, $1}' | xargs); do
    instanceWorks=0
    echo "Looking for instances for ${lb}"

    for instance in $(gcloud compute target-pools describe "$lb" | bosh int - --path '/instances?' | awk '{print $2}' | xargs); do
        echo "Looking for instance $instance"
        if gcloud compute instances describe "$instance" > /dev/null 2> /dev/null; then
            echo "Instance found [$instance] for loadbalancer [$lb]"
            instanceWorks=1
        fi
    done

    if [[ "$instanceWorks" == "0" ]]; then
        forwardingRules="${forwardingRules} ${lb/targetPools/forwardingRules}"
        targetPools="${targetPools} ${lb}"
    fi
done

if [[ "$forwardingRules" != ""  ]]; then
    # not all target-pools have fowarding-rules associated, so don't fail if it doesn't exist
    gcloud -q compute forwarding-rules delete $forwardingRules || true
fi

if [[ "$targetPools" != "" ]]; then
    gcloud -q compute target-pools delete $targetPools
fi
