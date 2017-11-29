#!/bin/bash

set -euxo pipefail

delete_vms() {
  network_id=$(bosh-cli int "$ENV_FILE" --path='/net_id')
  network_name=$(openstack network list -f value | grep "$network_id" | awk '{print $2}')
  server_list_with_details=$(openstack server list -f value)
  # Don't fail if no servers are found
  server_names=$(echo "$server_list_with_details" | grep "$network_name" | awk '{print $1}')

  openstack volume list --status available -c ID -f value | awk '{print $1}' | xargs -I{} openstack volume delete {}

  internal_ip=$(bosh-cli int "$ENV_FILE" --path='/internal_ip')
  port_list_with_details=$(openstack port list)
  echo "$port_list_with_details" | grep "$internal_ip" | awk '{print $2}' | xargs -I{} openstack port delete {} || echo 'No director port found'

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
        openstack volume delete "${volume_name}"
        echo "The volume became available and was deleted"
      else
        echo "The volume never became available and wasn't deleted"
      fi
    fi

    openstack server delete "${server_name}"
  done
}


OS_REGION_NAME=$(bosh-cli int "$ENV_FILE" --path='/region')
OS_PROJECT_NAME=$(bosh-cli int "$ENV_FILE" --path='/openstack_project')
OS_PROJECT_ID=$(bosh-cli int "$ENV_FILE" --path='/openstack_project_id')
OS_PASSWORD=$(bosh-cli int "$ENV_FILE" --path='/openstack_password')
OS_AUTH_URL=$(bosh-cli int "$ENV_FILE" --path='/auth_url')
OS_USERNAME=$(bosh-cli int "$ENV_FILE" --path='/openstack_username')
OS_USER_DOMAIN_NAME=$(bosh-cli int "$ENV_FILE" --path='/openstack_domain')
OS_IDENTITY_API_VERSION=3

export OS_REGION_NAME OS_PROJECT_NAME OS_PROJECT_ID OS_PASSWORD OS_AUTH_URL OS_USERNAME OS_USER_DOMAIN_NAME OS_IDENTITY_API_VERSION

delete_vms
