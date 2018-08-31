#! /usr/bin/env ruby

ops_files = '-o git-kubo-deployment/manifests/ops-files/use-runtime-config-bosh-dns.yml\
 -o git-kubo-deployment/manifests/ops-files/rename.yml\
 -o git-kubo-deployment/manifests/ops-files/vm-types.yml\
 -o git-kubo-ci/manifests/ops-files/add-api-server-endpoint.yml\
 -o git-kubo-deployment/manifests/ops-files/addons-spec.yml '
vars_files = '-l gcs-load-balancer-vars/load-balancer-vars.yml -l kubo-lock/metadata '
vars = "-v deployment_name=#{ENV['DEPLOYMENT_NAME']} -v worker_vm_type=worker -v master_vm_type=master "

var_file = "--var-file=addons-spec=#{ENV['ADDONS_SPEC']}"

if ENV['ENABLE_MULTI_AZ_TESTS']
  ops_files << '-o git-kubo-ci/manifests/ops-files/enable-multiaz-workers.yml '
else
  ops_files << '-o git-kubo-deployment/manifests/ops-files/misc/single-master.yml '
end

if ENV['IAAS'] =~ /^gcp/
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/gcp/cloud-provider.yml '
end

if ENV['IAAS'] =~ /^vsphere/
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/vsphere/cloud-provider.yml '
  vars_files << '-l director_uuid/var.yml '
end


if ENV['IAAS'] =~ /^vsphere-proxy/
  ops_files << '-o git-kubo-deployment/manifests/ops-files/add-proxy.yml '
end

if ENV['IAAS'] =~ /^aws/
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/aws/cloud-provider.yml '
  ops_files << '-o git-kubo-deployment/manifests/ops-files/iaas/aws/lb.yml '
end

cmd = ['bosh -n -d',
       ENV['DEPLOYMENT_NAME'],
       'deploy',
       '--no-redact',
       ENV['CFCR_MANIFEST_PATH'],
       ops_files,
       vars_files,
       vars,
       var_file].join(' ')
puts "command: #{cmd}"
File.write(ENV['BOSH_DEPLOY_COMMAND'], "#!/usr/bin/env bash\n" + cmd)
system("chmod +x #{ENV['BOSH_DEPLOY_COMMAND']}")
