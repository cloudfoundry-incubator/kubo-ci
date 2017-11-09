package cloudfoundry_test

import (
	"tests/test_helpers"
	"strings"
	"os"
	"github.com/onsi/gomega/gexec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	authenticationPolicy    string
	runner                  *test_helpers.KubectlRunner
}

const (
	abac = "ABAC"
	rbac = "RBAC"
)

func InitializeTestConfig(runner *test_helpers.KubectlRunner) IngressTestConfig {
	tc := IngressTestConfig{}
	tc.ingressSpec = test_helpers.PathFromRoot("specs/ingress.yml")
	tc.ingressRoles = test_helpers.PathFromRoot("specs/ingress-rbac-roles.yml")
	tc.rbacIngressSpec = test_helpers.PathFromRoot("specs/ingress-rbac.yml")
	tc.rbacServiceAccount = "nginx-ingress-serviceaccount"
	tc.runner = runner

	tc.tcpPort = os.Getenv("INGRESS_CONTROLLER_TCP_PORT")
	if tc.tcpPort == "" {
		Fail("Correct INGRESS_CONTROLLER_TCP_PORT has to be set")
	}
	tc.kubernetesServiceHost = os.Getenv("KUBERNETES_SERVICE_HOST")
	if tc.kubernetesServiceHost == "" {
		Fail("Correct KUBERNETES_SERVICE_HOST has to be set")
	}
	tc.kubernetesServicePort = os.Getenv("KUBERNETES_SERVICE_PORT")
	if tc.kubernetesServicePort == "" {
		Fail("Correct KUBERNETES_SERVICE_PORT has to be set")
	}
	tc.tlsKubernetesCert = os.Getenv("TLS_KUBERNETES_CERT")
	if tc.tlsKubernetesCert == "" {
		Fail("Correct TLS_KUBERNETES_CERT has to be set")
	}
	tc.tlsKubernetesPrivateKey = os.Getenv("TLS_KUBERNETES_PRIVATE_KEY")
	if tc.tlsKubernetesPrivateKey == "" {
		Fail("Correct TLS_KUBERNETES_PRIVATE_KEY has to be set")
	}
	tc.authenticationPolicy = strings.ToUpper(os.Getenv("KUBERNETES_AUTHENTICATION_POLICY"))
	if tc.authenticationPolicy != rbac && tc.authenticationPolicy != abac {
		Fail("Correct KUBERNETES_AUTHENTICATION_POLICY has to be set")
	}

	return tc
}

func (tc IngressTestConfig) createIngressController() {
	if tc.authenticationPolicy == rbac {
		tc.runner.RunKubectlCommandWithTimeout("create", "serviceaccount", tc.rbacServiceAccount)
		tc.runner.RunKubectlCommandWithTimeout("apply", "-f", tc.ingressRoles)
		tc.runner.RunKubectlCommandWithTimeout("create", "clusterrolebinding", "nginx-ingress-clusterrole-binding", "--clusterrole", "nginx-ingress-clusterrole", "--serviceaccount", tc.runner.Namespace()+":"+tc.rbacServiceAccount)
		tc.runner.RunKubectlCommandWithTimeout("create", "rolebinding", "nginx-ingress-role-binding", "--role", "nginx-ingress-role", "--serviceaccount", tc.runner.Namespace()+":"+tc.rbacServiceAccount)
		tc.runner.RunKubectlCommandWithTimeout("create", "-f", tc.rbacIngressSpec)
	} else {
		Eventually(tc.runner.RunKubectlCommand("create", "-f", tc.ingressSpec), "60s").Should(gexec.Exit(0))
	}
}

func (tc IngressTestConfig) deleteIngressController() {
	if tc.authenticationPolicy == rbac {
		Eventually(tc.runner.RunKubectlCommand("delete", "-f", tc.ingressRoles), "10s").Should(gexec.Exit())
		Eventually(tc.runner.RunKubectlCommand("delete", "clusterrolebinding", "nginx-ingress-clusterrole-binding"), "10s").Should(gexec.Exit())
	}
}
