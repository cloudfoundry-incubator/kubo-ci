package test_helpers

import (
	"fmt"
	"strings"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // load oidc auth
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ReadKubeConfig() (*restclient.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
}

func NewKubeClient() (k8s.Interface, error) {
	config, err := ReadKubeConfig()
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

func BearerToken() (string, error) {
	config, err := ReadKubeConfig()

	if err != nil {
		return "", err
	}

	token := ""

	if config.AuthProvider != nil && config.AuthProvider.Name == "oidc" {
		token = config.AuthProvider.Config["id-token"]
	} else {
		token = config.BearerToken
	}

	if token == "" {
		return "", fmt.Errorf("Token is empty")
	}
	return token, nil
}

func GetNodeIP() (string, error) {
	kubeclient, err := NewKubeClient()
	if err != nil {
		return "", err
	}

	nodes, err := kubeclient.CoreV1().Nodes().List(meta_v1.ListOptions{})
	if err != nil {
		return "", err
	}
	return nodes.Items[0].ObjectMeta.Labels["spec.ip"], nil

}
