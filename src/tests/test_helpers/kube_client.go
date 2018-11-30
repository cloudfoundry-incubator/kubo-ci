package test_helpers

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
	corev1 "k8s.io/api/core/v1"
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
	nodes, err := GetNodes()
	if err != nil {
		return "", err
	}

	providerID := nodes.Items[0].Spec.ProviderID
	iaas := strings.Split(providerID, ":")[0]
	switch iaas {
	case "vsphere", "gce", "openstack", "aws", "azure":
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
	nodes, err := GetNodes()
	if err != nil {
		return "", err
	}
	return nodes.Items[0].ObjectMeta.Labels["spec.ip"], nil
}

func GetNodes() (*corev1.NodeList, error) {
	kubeclient, err := NewKubeClient()
	if err != nil {
		return nil, err
	}

	return kubeclient.CoreV1().Nodes().List(meta_v1.ListOptions{})
}

func GetReadyNodes() ([]string, error) {
	nodes, err := GetNodes()
	if err != nil {
		return nil, err
	}
	readyNodes := []string{}

	for _, n := range nodes.Items {
		for _, condition := range n.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
				readyNodes = append(readyNodes, n.Name)
				break
			}
		}
	}

	return readyNodes, nil
}

func CreateTestNamespace(k8s k8s.Interface, prefix string) (*corev1.Namespace, error) {
	name := prefix + "-" + uuid.NewV4().String()
	labels := make(map[string]string)
	labels["test"] = prefix
	namespaceObject := corev1.Namespace{ObjectMeta: meta_v1.ObjectMeta{Name: name, Labels: labels}}
	return k8s.CoreV1().Namespaces().Create(&namespaceObject)
}

func DeleteTestNamespace(k8s k8s.Interface, namespace string) error {
	return k8s.CoreV1().Namespaces().Delete(namespace, &meta_v1.DeleteOptions{})
}

func HasWindowsWorkers() (bool, error) {
	nodes, err := GetNodes()
	if err != nil {
		return false, err
	}

	for _, n := range nodes.Items {
		if n.Status.NodeInfo.OperatingSystem == "windows" {
			return true, nil
		}
	}
	return false, nil
}
