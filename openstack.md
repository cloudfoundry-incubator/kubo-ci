# OpenStack guide

## Our OpenStack

The OpenStack installation the CFCR team uses is managed by the toolsmiths team. It is accessed through http://openstack-mitaka-02.mos9.cf-app.com and the login credentials is in LastPass note `kubo-openstack-ci`.

On the left hand nav, go to Identity -> Projects to see the list of all projects. The CFCR team uses the `clay` project. Projects sit above VMs in the organizational hierarchy. They allow operators to allow access to users access to OpenStack resources, and set quotas for the maximum number of volumes, security groups, floating IPs, etc.

Go to Project -> Compute -> Instances to see a list of all VMs in `clay`.

## OpenStack CLI

```bash
$ pip install python-openstackclient

$ export OS_AUTH_URL=http://openstack-mitaka-02.mos9.cf-app.com:5000/v2.0
$ export OS_PROJECT_NAME=clay
$ export OS_USERNAME=clay
$ export OS_PASSWORD=<password> # optional

# list networks
$ openstack network list

# list instances within the current project
$ openstack server list
```

> The OpenStack docs [lists every command](https://docs.openstack.org/python-openstackclient/latest/cli/command-list.html).

## Install PCF

1. Follow the guide on [Installing PCF on OpenStack](https://docs.pivotal.io/pivotalcf/1-12/customizing/openstack.html) on the Pivotal docs website.
1. Follow [Jaime's docs](https://docs.google.com/document/d/1PCnr4Lf0Y09OhW0yzPerorNrMPZQ7mAzA8vZNPd0oRU/edit#) on deploying CFCR on OpenStack.

## Deploy Concourse Worker

1. Log into the [OpenStack dashboard](http://openstack-mitaka-02.mos9.cf-app.com) as an admin.
1. On the left-hand navigation bar, click **Project** &rarr; **Network** &rarr; **Networks**.
1. Click the **+ Create Network** button on the top right corner.
1. Create a network for Concourse (i.e. `clay_net`)
1. Create a subnet with the CIDR block 192.168.130.0/24 and name it `concourse_subnet`.
1. Attach the network the router by clicking on the **Routers** page from the left navigation bar.
1. Click **+ Add Interface** and select the newly created subnet to the router.
1. Create a security group called `clay` and allow TCP access from everywhere.
1. Use scripts in `kubo-ci/concourse/scripts` to install the concourse worker (make sure the security group and network you created are being used).

## Creating a New Environment and Adding it to a Concourse Pool

1. Create a network and subnet - Go to Project -> Network -> Networks and click `Create Network`. Use the following values:
    * Network Name:  `example`
    * Subnet Name:  `example-subnet`
    * Network Address: `192.168.145.0/24` (Use a unique address if this one is taken)
    * Gateway IP: `192.168.145.1` (Must be in your subnet)
    * Enable DHCP: on
    * Allocation Pools: `192.168.145.2,192.168.145.254`
    * DNS Name Servers: `8.8.8.8`
    * Host Routes: none
1. Connect the network to a router - Go to Network -> Routers and click on `clay-router`. Click `Add Interface`. Use `example-subnet` as the Subnet and leave everything as default.
1. Add a lock file to the [`kubo-locks`](https://github.com/pivotal-cf-experimental/kubo-locks) repository. Navigate to the directory `kubo-openstack/unclaimed`. Create a new lock file by copy-and-pasting another lock file from this directory. Name your new lock file after the environment, in our case it would be `example`. Here're the properties you'll need to update:
    * `net_id`: `f351f900-16d0-426c-9616-c20e93c17e93`
    * `internal_ip`: `192.168.145.3`
    * `reserved_ips`: `192.168.145.1-192.168.145.2`
    * `director_name`: `example`
    * `internal_cidr`: `192.168.145.0/24`
    * `internal_gw`: `192.168.145.1`
    * `kubernetes_master_port`: [choose a unique port separate from the other locks]
1. Create an ops file for this environment in the [`kubo-odb-ci`](https://github.com/pivotal-cf-experimental/kubo-odb-ci) repository. Navigate to the directory `environments` and create a new directory and within that create a new ops file. In our case we would create a directory named `openstack-example`, and within we would create `openstack-example.yml`. Start by copy-and-pasting an ops file from another OpenStack environment. Here're the properties you'll need to update:
    * `.../service_catalog/id`: [generate a new GUID]
    * `.../plans/name=demo/plan_id`: [generate a new GUID]
    * `.../properties/broker_uri`: `http://openstack-example-odb.((cf_sys_domain))`
    * `.../routes/name=broker/uris`: `[ openstack-example-odb.((cf_sys_domain)) ]`
    * `.../nats/machines`: [leave the same]
