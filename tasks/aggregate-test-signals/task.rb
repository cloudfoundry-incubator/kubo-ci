#!/usr/bin/env ruby
# frozen_string_literal: true

require 'terminal-table'
require 'rainbow'
require 'set'

linux_files = Dir.glob('gcs-*-shipables').flat_map { |d| Dir.glob(d + '/*shipable') }
linux_files.delete_if { |f| /windows/.match(f) }
windows_files = Dir.glob('gcs-*windows*-shipables').flat_map { |d| Dir.glob(d + '/*shipable') }

def print_version_matrix(files)
  passed_matrix = {}
  version_overlap = File.read(files.first).split("\n")
  highest_versions = Set[]
  files.each do |f|
    version_overlap, passed_versions = find_overlap(f, version_overlap)
    build_name = f.split("/")[1].split("-shipable")[0]
    highest_version = version_overlap.last
    passed_matrix[build_name] = passed_versions
    highest_versions.add(highest_version)
  end

  rows = []
  passed_matrix.keys.each do |build|
    row = [build]
    highest_versions.each do |highest_version|
      if passed_matrix[build].include? highest_version
        row << "X"
      else
        row << ""
      end
    end
    rows << row
  end

  headers = ["build"]
  highest_versions.each do |highest_version|
    headers << highest_version
  end

  table = Terminal::Table.new :headings => headers, :rows => rows
  puts table
  puts

  version_overlap
end

def find_overlap(file, overlap)
  passed_versions = File.read(file).split("\n")
  overlap = passed_versions & overlap
  puts "After checking #{file} good versions are: #{overlap.last}"
  [ overlap, passed_versions ]
end

puts
puts "Looking for highest common green build..."
puts "....linux"
linux_overlap = print_version_matrix(linux_files)

if linux_overlap.any?
  linux_release_sha, deployment_sha, linux_build_number = linux_overlap.last.split
  File.write(ENV['SLACK_MESSAGE_FILE'],
             "Ready to :ship: <https://github.com/cloudfoundry-incubator/kubo-release/tree/#{linux_release_sha}/|#{linux_release_sha}> <https://github.com/cloudfoundry-incubator/kubo-deployment/tree/#{deployment_sha}/|#{deployment_sha}> Build number is #{linux_build_number}")
else
  puts 'No good versions yet'
  File.write(ENV['SLACK_MESSAGE_FILE'], 'No shippable version found')
  exit 1
end

puts Rainbow("Good versions are: #{linux_release_sha}, #{deployment_sha}. Build number is #{linux_build_number}").green

puts
puts "....windows"
windows_overlap = print_version_matrix(windows_files)

if windows_overlap.any?
  windows_release_sha, windows_deployment_sha, windows_build_number = windows_overlap.last.split
end

puts Rainbow("Good versions are: #{windows_release_sha}, #{windows_deployment_sha}. Build number is #{windows_build_number}").green

puts
puts "Highest green build for each pipeline..."
rows = []
files = linux_files + windows_files
files.each do |f|
  builds = File.read(f).split("\n")
  pipeline = f.split("/")[1].split("-shipable")[0]
  kubo_release_sha, kubo_deployment_sha, table_build_number = builds.last.split
  rows << [pipeline, kubo_release_sha, kubo_deployment_sha, table_build_number]
end
table = Terminal::Table.new :headings => ['Pipeline', 'Kubo Release Sha (Linux or Windows)', 'Kubo Deployment Sha', 'Build Number'], :rows => rows
puts table
puts

rows = []
rows << ["linux release sha", linux_release_sha]
rows << ["windows release sha", windows_release_sha]
rows << ["deployment sha", deployment_sha]
rows << ["linux build number", linux_build_number]
rows << ["windows build number", windows_build_number]
table = Terminal::Table.new :headings => ['Key', 'Value'], :rows => rows
puts table
puts

if linux_build_number == windows_build_number
  File.write(ENV['SHIPABLE_VERSION_FILE'], [linux_release_sha, windows_release_sha, deployment_sha, linux_build_number])
  puts "#{ENV['SHIPABLE_VERSION_FILE']} contains..."
  puts File.read(ENV['SHIPABLE_VERSION_FILE'])
else
  puts Rainbow("linux_build_number #{linux_build_number} does not match the windows_build_number #{windows_build_number} ").red
  puts Rainbow("would have shipped: #{linux_release_sha}, #{windows_release_sha}, #{deployment_sha}, #{linux_build_number}]")
end
