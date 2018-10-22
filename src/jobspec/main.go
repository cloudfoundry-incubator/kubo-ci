package main

import (
	"fmt"

	"k8s.io/kubernetes/jobspec/flag_generator"

	apiserver "k8s.io/kubernetes/cmd/kube-apiserver/app/options"

	"os"

	kubelet "k8s.io/kubernetes/cmd/kubelet/app/options"
)

func main() {
	component := os.Args[1]
	specPath := os.Args[2]
	jobSpec, _ := flag_generator.ReadSpecFile(specPath)

	switch component {
	case "kube-apiserver":
		apiserverBlacklistedFlags := []string{
			"apiserver-count",
			"cloud-provider",
			"cloud-config",
		}
		apiserverFlags := apiserver.NewServerRunOptions()
		jobSpec.Properties["k8s-args"] = flag_generator.GenerateArgsFromFlags(apiserverFlags, apiserverBlacklistedFlags)
	case "kubelet":
		kubeletBlacklistedFlags := []string{
			"cloud-config",
			"cloud-provider",
			"hostname-override",
			"config",
			"node-labels",
		}
		kubeletFlags, _ := kubelet.NewKubeletServer()
		jobSpec.Properties["k8s-args"] = flag_generator.GenerateArgsFromFlags(kubeletFlags, kubeletBlacklistedFlags)
	default:
		fmt.Printf("Unknown/unsupported kubernetes component: %s\n", component)
		os.Exit(1)
	}

	flag_generator.WriteSpecFile(specPath, jobSpec)
}
