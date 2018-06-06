package test_helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubeClient() (k8s.Interface, error) {
	var kubeconfig string
	home := os.Getenv("HOME")
	kubeconfig = filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return k8s.NewForConfig(config)
}

func IaaS() (string, error) {
	kubeclient, err := NewKubeClient()
	if err != nil {
		return "", err
	}

	nodes, err := kubeclient.CoreV1().Nodes().List(meta_v1.ListOptions{})
	if err != nil {
		return "", err
	}

	providerID := nodes.Items[0].Spec.ProviderID
	iaas := strings.Split(providerID, ":")[0]
	switch iaas {
	case "vsphere", "gce", "openstack", "aws":
		return iaas, nil
	}
	return "", fmt.Errorf("'%s' is not a valid iaas", iaas)
}
