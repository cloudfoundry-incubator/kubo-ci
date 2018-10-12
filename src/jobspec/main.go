package main

import (
	"k8s.io/kubernetes/jobspec/flag_generator"

	"k8s.io/kubernetes/cmd/kube-apiserver/app/options"

	"os"
)

func main() {
	blacklistedFlags := []string{
		"apiserver-count",
		"cloud-provider",
		"cloud-config",
	}
	specPath := os.Args[1]
	jobSpec := flag_generator.ReadExistingSpec(specPath)

	apiserverFlags := options.NewServerRunOptions()
	jobSpec.Properties["k8s-args"] = flag_generator.GenerateArgsFromFlags(apiserverFlags, blacklistedFlags)

	flag_generator.WriteNewSpec(specPath, jobSpec)
}
