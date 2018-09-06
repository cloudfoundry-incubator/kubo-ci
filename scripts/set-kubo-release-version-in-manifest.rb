require 'yaml'
require 'digest'
require 'fileutils'

def update_manifest(manifest)
  parsed_manifest = YAML.safe_load(manifest)
  parsed_manifest['releases'].delete_if { |r| r['name'] == 'kubo' }
  if ENV['IS_FINAL']
    parsed_manifest['releases'] << kubo_block
  else
    parsed_manifest['releases'] << kubo_dev_block
  end
  parsed_manifest
end

def kubo_block
  ver = version
  {
    'name' => 'kubo',
    'version' => ver,
    'sha1' => sha(ver),
    'url' => url(ver)
  }
end

def kubo_dev_block
  {
    'name' => 'kubo',
    'version' => 'latest'
  }
end

def version
  File.read('kubo-version/version')
end

def sha(version)
  shasum = Digest::SHA1.file "kubo-release-tarball/kubo-release-#{version}.tgz"
  shasum.hexdigest
end

def url(version)
  "https://github.com/cloudfoundry-incubator/kubo-release/releases/download/v#{version}/kubo-release-#{version}.tgz"
end

if $0 == __FILE__
  FileUtils.copy_entry 'git-kubo-deployment', 'git-kubo-deployment-output'
  File.open('git-kubo-deployment-output/manifests/cfcr.yml', 'w+') do |file|
    updated_manifest = update_manifest(file)
    file.write(updated_manifest.to_yaml)
  end
end
