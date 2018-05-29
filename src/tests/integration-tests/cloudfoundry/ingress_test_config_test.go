package cloudfoundry_test

import (
	"strconv"
	"tests/config"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type IngressTestConfig struct {
	tcpPort                 string
	tlsKubernetesCert       string
	tlsKubernetesPrivateKey string
	kubernetesServiceHost   string
	kubernetesServicePort   string
	ingressSpec             string
	ingressRoles            string
	rbacIngressSpec         string
	rbacServiceAccount      string
	runner                  *test_helpers.KubectlRunner
}

func InitializeIngressTestConfig(runner *test_helpers.KubectlRunner, testconfig config.Kubernetes) IngressTestConfig {
	tc := IngressTestConfig{}
	tc.ingressSpec = test_helpers.PathFromRoot("specs/ingress.yml")
	tc.ingressRoles = test_helpers.PathFromRoot("specs/ingress-rbac-roles.yml")
	tc.rbacIngressSpec = test_helpers.PathFromRoot("specs/ingress-rbac.yml")
	tc.rbacServiceAccount = "nginx-ingress-serviceaccount"
	tc.runner = runner

	tc.tcpPort = strconv.Itoa(testconfig.MasterPort + 20)
	tc.kubernetesServiceHost = testconfig.MasterHost
	if tc.kubernetesServiceHost == "" {
		Fail("Correct Kubernetes Master Host must be set in test config")
	}
	tc.kubernetesServicePort = strconv.Itoa(testconfig.MasterPort)
	if tc.kubernetesServicePort == "" {
		Fail("Correct Kubernetes Master Port must be set in test config")
	}
	tc.tlsKubernetesCert = testconfig.TLSCert
	if tc.tlsKubernetesCert == "" {
		Fail("Correct Kubernetes TLS Certificate must be set in test config")
	}
	tc.tlsKubernetesPrivateKey = testconfig.TLSPrivateKey
	if tc.tlsKubernetesPrivateKey == "" {
		Fail("Correct Kubernetes TLS Private Key must be set in test config")
	}
	return tc
}

func (tc IngressTestConfig) createIngressController() {
	tc.runner.RunKubectlCommandWithTimeout("create", "serviceaccount", tc.rbacServiceAccount)
	tc.runner.RunKubectlCommandWithTimeout("apply", "-f", tc.ingressRoles)
	tc.runner.RunKubectlCommandWithTimeout("create", "clusterrolebinding", "nginx-ingress-clusterrole-binding", "--clusterrole", "nginx-ingress-clusterrole", "--serviceaccount", tc.runner.Namespace()+":"+tc.rbacServiceAccount)
	tc.runner.RunKubectlCommandWithTimeout("create", "rolebinding", "nginx-ingress-role-binding", "--role", "nginx-ingress-role", "--serviceaccount", tc.runner.Namespace()+":"+tc.rbacServiceAccount)
	tc.runner.RunKubectlCommandWithTimeout("create", "-f", tc.rbacIngressSpec)
}

func (tc IngressTestConfig) deleteIngressController() {
	Eventually(tc.runner.RunKubectlCommand("delete", "-f", tc.ingressRoles), "10s").Should(gexec.Exit())
	Eventually(tc.runner.RunKubectlCommand("delete", "clusterrolebinding", "nginx-ingress-clusterrole-binding"), "10s").Should(gexec.Exit())
}
