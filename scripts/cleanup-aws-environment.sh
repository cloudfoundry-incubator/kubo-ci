#!/bin/bash

set -eu -o pipefail

lock_file="kubo-lock-repo/${POOL_NAME}/claimed/${ENV_NAME}"
director_ip=$(bosh-cli int ${lock_file} --path=/internal_ip)
subnet_id=$(bosh-cli int ${lock_file} --path=/subnet_id)
access_key_id=$(bosh-cli int ${lock_file} --path=/access_key_id)
secret_access_key=$(bosh-cli int ${lock_file} --path=/secret_access_key)
region=$(bosh-cli int ${lock_file} --path=/region)

mkdir -p ~/.aws

cat > ~/.aws/credentials <<-EOF
[default]
aws_access_key_id=${access_key_id}
aws_secret_access_key=${secret_access_key}
EOF

cat > ~/.aws/config <<-EOF
[default]
region=${region}
output=text
EOF

director_instance_id=$(aws ec2 describe-instances --query 'Reservations[*].Instances[*].InstanceId' --output text --filters "Name=network-interface.addresses.private-ip-address,Values=${director_ip}")
if [ -z "$director_instance_id" ]; then
  echo "No instance found for BOSH Director IP address"
else
  aws ec2 terminate-instances --instance-ids "$director_instance_id"
  aws ec2 wait instance-terminated --filters "Name=instance-id,Values=${director_instance_id}"
fi

instance_ids=$(aws ec2 describe-instances --query 'Reservations[*].Instances[*].InstanceId' --filters "Name=subnet-id,Values=${subnet_id}")
if [ -z "$instance_ids" ]; then
  echo "No instances found in subnet '${subnet_id}'"
else
  aws ec2 terminate-instances --instance-ids $instance_ids
  aws ec2 wait instance-terminated --filters "Name=instance-ids,Values=${instance_ids}"
fi
