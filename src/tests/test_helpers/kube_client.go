package test_helpers

import (
	"os"
	"path/filepath"

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
