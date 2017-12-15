# OpenStack guide

## Our OpenStack

The OpenStack installation the CFCR team uses is managed by the toolsmiths team. It is accessed through https://openstack-01.pez.pivotal.io and the login credentials is in LastPass note `Openstack Pez Dashboard user account (pcf-kubo)`.

Go to Identity -> Projects to see the list of all projects. The CFCR team uses `pcf-kubo-project`. Projects sit above VMs in the organizational hierarchy. They allow operators to allow access to users access to OpenStack resources, and set quotas for the maximum number of volumes, security groups, floating IPs, etc.

Go to Project -> Compute -> Instances to see a list of all VMs in the project.

## OpenStack CLI

From the OpenStack dashboard, go to `Project > Compute > Access & Security > API Access tab` and download the RC file (v3).

```bash
$ pip install python-openstackclient

$ source openrc.sh
# provide the password for pcf-kubo service account

# list networks
$ openstack network list

# list instances within the current project
$ openstack server list
```

> The OpenStack docs [lists every command](https://docs.openstack.org/python-openstackclient/latest/cli/command-list.html).

## Install PCF

1. Follow the guide on [Installing PCF on OpenStack](https://docs.pivotal.io/pivotalcf/1-12/customizing/openstack.html) on the Pivotal docs website.
   - Ensure DHCP is enabled when creating your subnet
   - Ensure all subnets have the DNS `8.8.8.8`
   - To create a floating IP, you must first create an instance.  The subnet should be connected to the router (create an interface from the router to the subnet 'port').  From the instance UI, you can then associate a floating IP.
   - The validator step requires an image to be uploaded. You'll have to download the .img file for [Ubuntu](https://docs.openstack.org/image-guide/obtain-images.html) and then create the image from the command line (it does not work from the Web GUI) `openstack image create --file ~/Downloads/image-file-name.img --disk-format raw "ubuntu-xenial"`.
   - Creating an image in step 4 may require you to do this from bash: `openstack image create --file ~/Downloads/pcf-openstack-1.12.5.raw --disk-format raw --private --protected --min-disk 20 --min-ram 8192 "ops manager 1.12.5"`
   - To set up a FQDN (e.g. `openstack-pez-01.cf-app.com`) for Ops Manager in the external DNS, you will need the `Dev DNS Management (AWS)` credentials in LastPass, in the `Shared-Opensource Common` folder.  Log in to AWS and go to Route53.
1. Follow the guide on [Installing Elastic Runtime](https://docs.pivotal.io/pivotalcf/1-12/customizing/openstack-er-config.html) on the Pivotal docs website.
  - Create two new floating IP addresses (Compute > Access & Security > Floating IPs).  Go to Route 53 (use the `Dev DNS Management (AWS)` credentials in LastPass, in the `Shared-Opensource Common` folder) and create a wildcard domain for systems and applications; e.g. `*.sys.openstack-pez-01.cf-app.com` and `*.app.openstack-pez-01.cf-app.com` assigning both to one of the floating IP addresses.  Then create `tcp.openstack-pez-01.cf-app.com` and assign that to the other IP address.
  - Only the following tabs and parameters of Pivotal Elastic Runtime options need to be altered.  Retain the defaults unless noted otherwise below:
    - _Assign AZs and Networks_: save defaults
    - _Domains_: enter the wildcard domains created in Route53 (described above)
    - _Networking_:
      - _Certificate and Private Key for HAProxy and Router_: Click __Generate RSA Certificate__ and give it the root FQDN e.g. `*.openstack-pez-01.cf-app.com`
      - _TLS Cipher Suites for Router_: Copy the default from the tooltip shown when you click in the textbox
      - _TLS Cipher Suites for HAProxy_: Copy the default from the tooltip shown beneath the textbox
      - _HAProxy forwards requests to Router over TLS..._: Disable
      - _Disable SSL certificate verification for this environment_: Select this checkbox
      - _Enable TCP Routing_: Select this checkbox
      - _TCP Routing Ports_: `1025-4096`
    - _Application Security Groups_: type `x`
    - _UAA_: _SAML Service Provider Credentials_: Click __Generate RSA Certificate__ and give it the root FQDN e.g. `*.openstack-pez-01.cf-app.com`
    - _Internal MySQL_:
      - _Email address (required)_: type `pcf-kubo@pivotal.io`
      - _Server Activity Logging_: disable
    - _Resource Config_:
      - Add the floating IP that points to the `*.sys` & `*.app` to the _Router_
      - Add the floating IP that points to `tcp.` to the _TCP Router_
      
You may also wish to read [Jaime's docs](https://docs.google.com/document/d/1PCnr4Lf0Y09OhW0yzPerorNrMPZQ7mAzA8vZNPd0oRU/edit#) on deploying CFCR on OpenStack.

## Create a jumpbox

1. Locally create a ssh key pair using `ssh-keygen`, and save this to OpenStack, Compute>Access & Security>Key Pairs.
1. Store this key pair in Lastpass in the shared `jumpbox` folder
1. Create a new security group called `jumpbox`, and give it rules for SSH ingress only but all egress.
1. Create a new instance (use an ubuntu trusty image), assign it to the `jumpbox` security group and in the Key Pair dialog select the newly created key pair.
1. Run the `sync-jumpbox` script in kubo-home repo.  This requires you to be logged in to the Lastpass CLI.

## Deploy Concourse Worker

1. Log into the [OpenStack dashboard](https://openstack-01.pez.pivotal.io) as an admin.
1. On the left-hand navigation bar, click **Project** &rarr; **Network** &rarr; **Networks**.
1. Click the **+ Create Network** button on the top right corner.
1. Create a network for Concourse (i.e. `concourse`)
1. Create a subnet with the CIDR block `192.168.1.0/24` and name it `concourse-subnet`.
1. Attach the network the router by clicking on the **Routers** page from the left navigation bar.
1. Click **+ Add Interface** and select the newly created subnet to the router.
1. Create a security group called `concourse` and allow: 
    1. TCP access from everywhere.
    1. UDP access from itself
1. Edit `opsmanager` security group to allow:
    1. All TCP access from `concourse` security group
    1. All TCP access from Floating IP CIDR.
1. `sshuttle` into the jumpbox created in the section above.
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
1. [Create a Routing UAA client](https://docs-cfcr.cfapps.io/installing/cf-routing/#step-2-create-a-routing-uaa-client)
1. Add a lock file to the [`kubo-locks`](https://github.com/pivotal-cf-experimental/kubo-locks) repository. Navigate to the directory `kubo-openstack/unclaimed`. Create a new lock file by copy-and-pasting another lock file from this directory. Name your new lock file after the environment, in our case it would be `example`. Here're the properties you'll need to update:
    * `net_id`: `f351f900-16d0-426c-9616-c20e93c17e93` The ID of the network in which the environment will create VMs
    * `internal_ip`: `192.168.145.3` Ensure this is within the CIDR of the network specified below.  The master will be automatically created at this IP.
    * `reserved_ips`: `192.168.145.1-192.168.145.2` . Within the CIDR of the network specified below
    * `director_name`: `example` .  The name of your environment
    * `internal_cidr`: `192.168.145.0/24` The CIDR of the subnet related to the above network ID
    * `internal_gw`: `192.168.145.1` Default is `1`, within your CIDR
    * `kubernetes_master_port`: [choose a unique port separate from the other locks] You will have to create a floating IP
    * `private_key`: The private rsa key for the bosh director
    * For `routing_mode: cf` environments:
        * `kubernetes_master_host`: e.g. `tcp.openstack-pez-01.cf-app.com`. Use the FQDN you set up in Route53, described above.
        * `routing_cf_api_url`: e.g. `https://api.sys.openstack-pez-01.cf-app.com`. Prepend `https://api.` to the FQDN for sys you set up in Route53, described above.
        * `routing_cf_client_id`: `routing_api_client`
        * `routing-cf-client-secret`: [The passphrase set up in the previous step for the UAA client]
        * `routing_cf_uaa_url`: e.g. `https://uaa.sys.openstack-pez-01.cf-app.com`. Prepend `https://uaa.` to the FQDN for sys you set up in Route53, described above.
        * `routing_cf_app_domain_name`: e.g. `app.openstack-pez-01.cf-app.com`. Use the FQDN for apps you set up in Route53, described above
        * `routing-cf-sys-domain-name`: e.g. `sys.openstack-pez-01.cf-app.com`. Use the FQDN for sys you set up in Route53, described above
        * `routing_cf_nats_internal_ips`: This can be found in the Ops Manager Pivotal Elastic Runtime tile, in the Status tab. For a 'full' footprint deployment the IP is the IP for the NATS VM, for a small footprint deployment it is the [IP for the Database VM](https://docs.pivotal.io/pivotalcf/1-12/customizing/small-footprint.html).
        * `routing-cf-nats-password`: Can be found in Ops Manager tile for Pivotal Elastic Runtime, in the Credentials tab.  Go to Jobs>NATS>Credentials.
1. Create an ops file for this environment in the [`kubo-odb-ci`](https://github.com/pivotal-cf-experimental/kubo-odb-ci) repository. Navigate to the directory `environments` and create a new directory and within that create a new ops file. In our case we would create a directory named `openstack-example`, and within we would create `openstack-example.yml`. Start by copy-and-pasting an ops file from another OpenStack environment. Here're the properties you'll need to update:
    * `.../service_catalog/id`: [generate a new GUID]
    * `.../plans/name=demo/plan_id`: [generate a new GUID]
    * `.../properties/broker_uri`: `http://openstack-example-odb.((cf_sys_domain))`
    * `.../routes/name=broker/uris`: `[ openstack-example-odb.((cf_sys_domain)) ]`
    * `.../nats/machines`: [leave the same]
