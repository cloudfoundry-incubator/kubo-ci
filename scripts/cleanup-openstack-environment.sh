#!/bin/bash

set -xue

delete_by_ip() {
    local address=$1
    local server_name=$(openstack server list -f value --ip ${address} | awk '{print $1}')
    if openstack server show ${server_name} -f yaml | bosh-cli int --path=/properties - | grep -E "job='(master|worker|etcd|route-sync|bosh)"; then
      local volume_name=$(openstack server show "${server_name}" -f yaml | bosh-cli int --path /volumes_attached - | cut -d "'" -f2)
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
    fi

}
delete_vms() {
  local internal_cidr=$(bosh-cli int "$ENV_FILE" --path=/internal_cidr)
  local network_prefix=$(echo ${internal_cidr} | awk -F "." '{print $1"."$2"."$3"."}')
  local min_ip_char=$(ipcalc ${internal_cidr} | grep HostMin | awk '{print $2}' | awk -F "." '{print $4}')
  local max_ip_char=$(ipcalc ${internal_cidr} | grep HostMax | awk '{print $2}' | awk -F "." '{print $4}')

  openstack volume list --status available -c ID -f value | awk '{print $1}' | xargs -I{} openstack volume delete {}

  # Delete the director first
  internal_ip=$(bosh-cli int "$ENV_FILE" --path='/internal_ip')
  delete_by_ip "$internal_ip"

  for address in ${network_prefix}{${min_ip_char}..${max_ip_char}}; do
    delete_by_ip "$address"
  done
}


export OS_REGION_NAME
OS_REGION_NAME=$(bosh-cli int "$ENV_FILE" --path='/region')
export OS_PROJECT_NAME
OS_PROJECT_NAME=$(bosh-cli int "$ENV_FILE" --path='/openstack_project')
export OS_PASSWORD
OS_PASSWORD=$(bosh-cli int "$ENV_FILE" --path='/openstack_password')
export OS_AUTH_URL
OS_AUTH_URL=$(bosh-cli int "$ENV_FILE" --path='/auth_url')
export OS_USERNAME
OS_USERNAME=$(bosh-cli int "$ENV_FILE" --path='/openstack_username')

delete_vms
