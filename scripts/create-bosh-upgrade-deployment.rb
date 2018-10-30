#! /usr/bin/env ruby

ops_files = '-o git-kubo-deployment/manifests/ops-files/use-runtime-config-bosh-dns.yml\
 -o git-kubo-deployment/manifests/ops-files/rename.yml\
 -o git-kubo-deployment/manifests/ops-files/vm-types.yml\
 -o git-kubo-deployment/manifests/ops-files/misc/dev.yml \
 -o git-kubo-deployment/manifests/ops-files/enable-nfs.yml \
 -o git-kubo-ci/manifests/ops-files/add-api-server-endpoint.yml '
vars_files = '-l gcs-load-balancer-vars/load-balancer-vars.yml -l kubo-lock/metadata '
vars = "-v deployment_name=#{ENV['DEPLOYMENT_NAME']} -v worker_vm_type=worker -v master_vm_type=master -v apply_addons_vm_type=minimal"


if ENV['ENABLE_MULTI_AZ_TESTS'] != 'false'
  ops_files << '-o git-kubo-ci/manifests/ops-files/enable-multiaz-workers.yml '
else
  ops_files << '-o git-kubo-deployment/manifests/ops-files/misc/single-master.yml '
  ops_files << '-o git-kubo-ci/manifests/ops-files/scale-to-one-az.yml '
end

if ENV['IAAS'] =~ /^gcp/
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/gcp/cloud-provider.yml '
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/gcp/add-subnetwork-for-internal-load-balancer.yml '
  ops_files << '-o git-kubo-deployment/manifests/ops-files/use-vm-extensions.yml '
end

if ENV['IAAS'] =~ /^vsphere/
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/vsphere/cloud-provider.yml '
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/vsphere/use-vm-extensions.yml '
  vars_files << '-l director_uuid/var.yml '
end


if ENV['IAAS'] =~ /^vsphere-proxy/
  ops_files << '-o git-kubo-deployment/manifests/ops-files/add-proxy.yml '
  ops_files << '-o git-kubo-ci/manifests/ops-files/add-master-static-ips.yml '
end

if ENV['IAAS'] =~ /^aws/
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/aws/cloud-provider.yml '
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/aws/lb.yml '
  ops_files << '-o git-kubo-deployment/manifests/ops-files/use-vm-extensions.yml '
end

cmd = ['bosh -n -d',
       ENV['DEPLOYMENT_NAME'],
       'deploy',
       '--no-redact',
       ENV['CFCR_MANIFEST_PATH'],
       ops_files,
       vars_files,
       vars].join(' ')
puts "command: #{cmd}"
File.write(ENV['BOSH_DEPLOY_COMMAND'], "#!/usr/bin/env bash\n" + cmd)
system("chmod +x #{ENV['BOSH_DEPLOY_COMMAND']}")
