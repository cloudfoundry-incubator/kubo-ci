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
	jobSpec, _ := flag_generator.ReadSpecFile(specPath)

	apiserverFlags := options.NewServerRunOptions()
	jobSpec.Properties["k8s-args"] = flag_generator.GenerateArgsFromFlags(apiserverFlags, blacklistedFlags)

	flag_generator.WriteSpecFile(specPath, jobSpec)
}
