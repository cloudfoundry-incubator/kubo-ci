#!/bin/bash

set -xue

delete_vms() {
  network_id=$(bosh-cli int "$lock_file" --path='/net_id')
  network_name=$(openstack network list -f value | grep "$network_id" | awk '{print $2}')
  server_names=$(openstack server list -f value | grep "$network_name" | awk '{print $1}')

  openstack volume list --status available -c ID -f value | awk '{print $1}' | xargs -I{} openstack volume delete {}

  internal_ip=$(bosh-cli int $lock_file --path='/internal_ip')
  openstack port list | grep "$internal_ip" | awk '{print $2}' | xargs -I{} openstack port delete {}
  internal_ip=$(bosh-cli int $lock_file --path='/internal_ip')
  openstack port list | grep "$internal_ip" | awk '{print $2}' | xargs -I{} openstack port delete {}

  for server_name in ${server_names}
  do
    volume_name=$(openstack server show "${server_name}" -f yaml | bosh-cli int --path /volumes_attached - | cut -d "'" -f2)
    if [ ! "${volume_name}" == "" ]; then
      openstack server remove volume "${server_name}" "${volume_name}"

      if timeout 120 /bin/bash <<EOF
        until openstack volume show "${volume_name}" -f yaml | bosh-cli int --path /status - | grep "available"; do
          sleep 2
        done
EOF
      then
        openstack volume delete "${volume_name}" --force
        echo "The volume became available and was deleted"
      else
        echo "The volume never became available and wasn't deleted"
      fi
    fi

    openstack server delete "${server_name}"
  done


  internal_ip=$(bosh-cli int "$lock_file" --path='/internal_ip')
  openstack port list -f value | grep "$internal_ip" | awk '{print $1}' | xargs -I{} openstack port delete {}
}

export lock_file="kubo-lock-repo/${POOL_NAME}/claimed/${ENV_NAME}"

export OS_REGION_NAME
OS_REGION_NAME=$(bosh-cli int "$lock_file" --path='/region')
export OS_PROJECT_NAME
OS_PROJECT_NAME=$(bosh-cli int "$lock_file" --path='/openstack_project')
export OS_PASSWORD
OS_PASSWORD=$(bosh-cli int "$lock_file" --path='/openstack_password')
export OS_AUTH_URL
OS_AUTH_URL=$(bosh-cli int "$lock_file" --path='/auth_url')
export OS_USERNAME
OS_USERNAME=$(bosh-cli int "$lock_file" --path='/openstack_username')

delete_vms
