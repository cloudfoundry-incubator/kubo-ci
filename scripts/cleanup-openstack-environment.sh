#!/bin/bash

set -xu # not using -e or pipefail, since list/delete of resources can potentially failed

delete_vms() {
  network_id=$(bosh-cli int $lock_file --path='/net_id')
  network_name=$(openstack network list | grep "$network_id" | awk '{print $4}')
  openstack server list | grep "$network_name" | awk '{print $2}' | xargs -I{} openstack server delete {}

  openstack volume list --status available -c ID | grep -v ID | grep -v '\-----' | awk '{print $2}' | xargs -I{} openstack volume delete {}

  internal_ip=$(bosh-cli int $lock_file --path='/internal_ip')
  openstack port list | grep "$internal_ip" | awk '{print $2}' | xargs -I{} openstack port delete {}
}

export lock_file="kubo-lock-repo/${POOL_NAME}/claimed/${ENV_NAME}"

export OS_REGION_NAME
OS_REGION_NAME=$(bosh-cli int $lock_file --path='/region')
export OS_PROJECT_NAME
OS_PROJECT_NAME=$(bosh-cli int $lock_file --path='/openstack_project')
export OS_PASSWORD
OS_PASSWORD=$(bosh-cli int $lock_file --path='/openstack_password')
export OS_AUTH_URL
OS_AUTH_URL=$(bosh-cli int $lock_file --path='/auth_url')
export OS_USERNAME
OS_USERNAME=$(bosh-cli int $lock_file --path='/openstack_username')

delete_vms
