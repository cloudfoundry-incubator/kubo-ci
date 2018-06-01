package test_helpers

import (
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
	return strings.Split(providerID, ":")[0], nil
}
