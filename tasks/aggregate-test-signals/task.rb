#!/usr/bin/env ruby
# frozen_string_literal: true

files = Dir.glob('gcs-*-shipables').flat_map { |d| Dir.glob(d + '/*shipable') }

puts "Highest green build for each pipeline..."
rows = []
files.each do |f|
  builds = File.read(f).split("\n")
  pipeline = f.split("/")[1].split("-shipable")[0]
  release_sha, deployment_sha, build_number = builds.last.split
  rows << [pipeline, build_number]
end
table = Terminal::Table.new :headings => ['Pipeline', 'Build Number'], :rows => rows
puts table
puts

puts "Looking for highest common green build..."
overlap = File.read(files.first).split("\n")
files.each do |f|
  overlap = File.read(f).split("\n") & overlap
  puts "After checking #{f} good versions are: #{overlap.last}"
end

if overlap.any?
  release_sha, deployment_sha, build_number = overlap.last.split
  File.write(ENV['SLACK_MESSAGE_FILE'],
             "Ready to :ship: <https://github.com/cloudfoundry-incubator/kubo-release/tree/#{release_sha}/|#{release_sha}> <https://github.com/cloudfoundry-incubator/kubo-deployment/tree/#{deployment_sha}/|#{deployment_sha}> Build number is #{build_number}")
  File.write(ENV['SHIPABLE_VERSION_FILE'], overlap.last)

else
  puts 'No good versions yet'
  File.write(ENV['SLACK_MESSAGE_FILE'], 'No shippable version found')
  exit 1
end
puts "Good versions are: #{release_sha}, #{deployment_sha}. Build number is #{build_number}"
