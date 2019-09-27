def wait_for_apply_to_finish(crd_count)
  puts "Wait for the expected number of crds..."
  count = check_crd_count
  while count != crd_count
    sleep 1
    count = check_crd_count
  end
end

def check_crd_count()
  command = "kubectl get crds | grep 'istio.io' | wc -l"
  puts "#{command}"
  output = `#{command}`
  crds = output.to_i
  puts "#{crds} crds"
  crds
end

wait_for_apply_to_finish(ARGV[0].to_i)
