#!/bin/bash

[ -z "$DEBUG" ] || set -x
set -eu -o pipefail

director_name="$(bosh int "${ENV_FILE}" --path=/director_name)"
director_ip="$(bosh int "${ENV_FILE}" --path=/internal_ip)"
subnet_id="$(bosh int "${ENV_FILE}" --path=/subnet_id)"
access_key_id="$(bosh int "${ENV_FILE}" --path=/access_key_id)"
secret_access_key="$(bosh int "${ENV_FILE}" --path=/secret_access_key)"
region="$(bosh int "${ENV_FILE}" --path=/region)"

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

delete_volumes() {
  local volume_ids=$1
  for volume in $volume_ids; do
    aws ec2 delete-volume --volume-id "$volume"
    aws ec2 wait volume-deleted --volume-ids "$volume"
  done
}

cleanup_load_balancers() {
  local is_k8s_lb
  local lbs=$(aws elb describe-load-balancers --output json | jq '.LoadBalancerDescriptions | .[] | .LoadBalancerName | select(test("^[a-z0-9]+$"))' -r)
  for lb in $lbs; do
    is_k8s_lb=$(aws elb describe-tags --load-balancer-name "${lb}" --output json | jq '.TagDescriptions[0].Tags | .[] | select(.Key == "kubernetes.io/cluster/'"${director_name}"'").Value' -r)
    if [[ "owned" == "${is_k8s_lb}" ]]; then
      echo "Deleting load balancer ${lb}"
      aws elb delete-load-balancer --load-balancer-name "${lb}"
    fi
  done
}

director_instance_id=$(aws ec2 describe-instances --query 'Reservations[*].Instances[*].InstanceId' --output text --filters "Name=network-interface.addresses.private-ip-address,Values=${director_ip}" "Name=subnet-id,Values=${subnet_id}")
if [ -z "$director_instance_id" ]; then
  echo "No instance found for BOSH Director IP address"
else
  aws ec2 terminate-instances --instance-ids "$director_instance_id"
  aws ec2 wait instance-terminated --filters "Name=instance-id,Values=${director_instance_id}"
fi

instance_ids=$(aws ec2 describe-instances --query 'Reservations[*].Instances[*].InstanceId' --filters "Name=tag:KubernetesCluster,Values=${director_name}")
if [ -z "$instance_ids" ]; then
  echo "No instances found with tag KubernetesCluster:'${director_name}'"
else
  aws ec2 terminate-instances --instance-ids ${instance_ids}
  aws ec2 wait instance-terminated --instance-ids ${instance_ids}
fi

volume_ids=$(aws ec2 describe-volumes --output text --no-paginate --query 'Volumes[*].VolumeId' --filters "Name=tag:director,Values=${director_name}" "Name=status,Values=available")
director_volume_ids=$(aws ec2 describe-volumes --output text --no-paginate --query 'Volumes[*].VolumeId' --filters "Name=tag:director_name,Values=${director_name}" "Name=status,Values=available")

if [ ! -z "$director_volume_ids" ]; then
  echo "Deleting volumes associated to director '${director_name}'"
  delete_volumes "$director_volume_ids"
fi

if [ ! -z "$volume_ids" ]; then
  echo "Deleting volumes tagged with director name '${director_name}'"
  delete_volumes "$volume_ids"
fi

echo "Deleting load balancers"
cleanup_load_balancers
