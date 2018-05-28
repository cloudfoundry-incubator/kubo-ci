#!/bin/bash

state_dir="${HOME}/workspace/kubo-locks/kubo-gcp-lb-lifecycle/terraform/"

mkdir -p "${state_dir}"

terraform init
for lock in ${HOME}/workspace/kubo-locks/kubo-gcp-lb-lifecycle/claimed/*; do
    base_name=$(basename $lock)
    env_name=${base_name#"gcp-"}
    terraform apply --auto-approve --var "service_account_key_path=gcp.json" \
        --var "projectid=cf-pcf-kubo" --var "region=us-central1" \
        --var "prefix=${env_name}" --state="${state_dir}/${env_name}.tfstate" \
        "${HOME}/workspace/kubo-ci/terraform/"

    ip_address="$(terraform output --state="${state_dir}/${env_name}.tfstate" ip_address)"
    change_tmp_file="$(mktemp)"
    cat > ${change_tmp_file} <<EOF
{
    "Comment": "Update record to reflect new IP address of home router",
    "Changes": [
        {
            "Action": "UPSERT",
            "ResourceRecordSet": {
                "Name": "${env_name}-gcp-lb.kubo.sh.",
                "Type": "A",
                "TTL": 300,
                "ResourceRecords": [
                    {
                        "Value": "${ip_address}"
                    }
                ]
            }
        }
    ]
}
EOF
done
