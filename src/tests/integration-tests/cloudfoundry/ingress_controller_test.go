package cloudfoundry_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"tests/test_helpers"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Testing Ingress Controller", func() {

	var (
		tcpPort                 string
		tlsKubernetesCert       string
		tlsKubernetesPrivateKey string
		kubernetesServiceHost   string
		kubernetesServicePort   string

		ingressSpec = test_helpers.PathFromRoot("specs/ingress.yml")
		runner      *test_helpers.KubectlRunner
		hasPassed   bool

		ingressRoles       = test_helpers.PathFromRoot("specs/ingress-rbac-roles.yml")
		rbacIngressSpec    = test_helpers.PathFromRoot("specs/ingress-rbac.yml")
		rbacServiceAccount = "nginx-ingress-serviceaccount"
	)

	const (
		authPolicyAttribute = "ABAC"
		authPolicyRole      = "RBAC"
	)

	createAbacIngressController := func() {
		Eventually(runner.RunKubectlCommand(
			"create", "-f", ingressSpec), "60s").Should(gexec.Exit(0))
	}

	createRbacIngressController := func() {
		Eventually(runner.RunKubectlCommand("create", "serviceaccount", rbacServiceAccount)).Should(gexec.Exit(0))
		Eventually(runner.RunKubectlCommand("apply", "-f", ingressRoles)).Should(gexec.Exit(0))
		Eventually(runner.RunKubectlCommand("create", "clusterrolebinding", "nginx-ingress-clusterrole-binding", "--clusterrole", "nginx-ingress-clusterrole", "--serviceaccount", runner.Namespace()+":"+rbacServiceAccount)).
			Should(gexec.Exit(0))
		Eventually(runner.RunKubectlCommand("create", "rolebinding", "nginx-ingress-role-binding", "--role", "nginx-ingress-role", "--serviceaccount", runner.Namespace()+":"+rbacServiceAccount)).
			Should(gexec.Exit(0))
		Eventually(runner.RunKubectlCommand("create", "-f", rbacIngressSpec), "60s").Should(gexec.Exit(0))
	}

	deleteRbacIngressController := func() {
		// Delete ingress roles and clusterrolebinding - everything else should be deleted with the namespace
		Eventually(runner.RunKubectlCommand("delete", "-f", ingressRoles)).Should(gexec.Exit())
		Eventually(runner.RunKubectlCommand("delete", "clusterrolebinding", "nginx-ingress-clusterrole-binding")).Should(gexec.Exit())
	}

	BeforeEach(func() {
		tcpPort = os.Getenv("INGRESS_CONTROLLER_TCP_PORT")
		if tcpPort == "" {
			Fail("Correct INGRESS_CONTROLLER_TCP_PORT has to be set")
		}
		kubernetesServiceHost = os.Getenv("KUBERNETES_SERVICE_HOST")
		if kubernetesServiceHost == "" {
			Fail("Correct KUBERNETES_SERVICE_HOST has to be set")
		}
		kubernetesServicePort = os.Getenv("KUBERNETES_SERVICE_PORT")
		if kubernetesServicePort == "" {
			Fail("Correct KUBERNETES_SERVICE_PORT has to be set")
		}
		tlsKubernetesCert = os.Getenv("TLS_KUBERNETES_CERT")
		if tlsKubernetesCert == "" {
			Fail("Correct TLS_KUBERNETES_CERT has to be set")
		}
		tlsKubernetesPrivateKey = os.Getenv("TLS_KUBERNETES_PRIVATE_KEY")
		if tlsKubernetesPrivateKey == "" {
			Fail("Correct TLS_KUBERNETES_PRIVATE_KEY has to be set")
		}

		authenticationPolicy := strings.ToUpper(os.Getenv("KUBERNETES_AUTHENTICATION_POLICY"))

		if authenticationPolicy != authPolicyAttribute && authenticationPolicy != authPolicyRole {
			authenticationPolicy = authPolicyRole
		}

		certFile, _ := ioutil.TempFile(os.TempDir(), "cert")
		_, err := certFile.WriteString(tlsKubernetesCert)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(certFile.Name())

		keyFile, _ := ioutil.TempFile(os.TempDir(), "key")
		_, err = keyFile.WriteString(tlsKubernetesPrivateKey)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(keyFile.Name())

		runner = test_helpers.NewKubectlRunner()
		runner.RunKubectlCommand(
			"create", "namespace", runner.Namespace()).Wait("60s")

		Eventually(
			runner.RunKubectlCommand(
				"create", "secret", "tls", "tls-kubernetes",
				"--cert", certFile.Name(),
				"--key", keyFile.Name(),
			),
			"60s",
		).Should(gexec.Exit(0))

		Eventually(
			runner.RunKubectlCommand(
				"create", "secret", "generic", "kubernetes-service",
				fmt.Sprintf("--from-literal=host=%s", kubernetesServiceHost),
				fmt.Sprintf("--from-literal=port=%s", kubernetesServicePort),
			),
			"60s",
		).Should(gexec.Exit(0))

		if authenticationPolicy == authPolicyRole {
			createRbacIngressController()
		} else {
			createAbacIngressController()
		}

		Eventually(runner.RunKubectlCommand(
			"rollout", "status", "deployments/default-http-backend", "-w"), "300s").Should(gexec.Exit(0))
		Eventually(runner.RunKubectlCommand(
			"rollout", "status", "deployments/nginx-ingress-controller", "-w"), "300s").Should(gexec.Exit(0))
		Eventually(runner.RunKubectlCommand(
			"rollout", "status", "deployments/simple-http-server", "-w"), "300s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		if hasPassed {
			authenticationPolicy := strings.ToUpper(os.Getenv("KUBERNETES_AUTHENTICATION_POLICY"))

			if authenticationPolicy != authPolicyAttribute && authenticationPolicy != authPolicyRole {
				authenticationPolicy = authPolicyAttribute
			}

			if authenticationPolicy == authPolicyRole {
				deleteRbacIngressController()
			}

			Eventually(runner.RunKubectlCommand(
				"delete", "-f", ingressSpec), "60s").Should(gexec.Exit())

			Eventually(runner.RunKubectlCommand(
				"delete", "secret", "tls-kubernetes")).Should(gexec.Exit())

			Eventually(runner.RunKubectlCommand(
				"delete", "secret", "kubernetes-service")).Should(gexec.Exit())

			runner.RunKubectlCommand(
				"delete", "namespace", runner.Namespace()).Wait("60s")
		}
	})

	It("Allows routing via Ingress Controller", func() {
		serviceName := test_helpers.GenerateRandomName()
		appUrl := fmt.Sprintf("http://%s.%s", serviceName, appsDomain)
		httpClient := http.Client{
			Timeout: time.Duration(5 * time.Second),
		}

		By("exposing it via HTTP - " + serviceName)
		result, err := httpClient.Get(appUrl)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.StatusCode).To(Equal(404))

		httpLabel := fmt.Sprintf("http-route-sync=%s", serviceName)
		Eventually(runner.RunKubectlCommand("label", "services", "nginx-ingress-controller", httpLabel), "10s").Should(gexec.Exit(0))

		Eventually(func() int {
			result, err := httpClient.Get(appUrl + "/simple-http-server")
			if err != nil {
				fmt.Println(err)
				return -1
			}
			return result.StatusCode
		}, "120s", "5s").Should(Equal(200))

		result, err = httpClient.Get(appUrl)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.StatusCode).To(Equal(404))

		By("exposing it via TCP")
		appUrl = fmt.Sprintf("http://%s:%s", tcpRouterDNSName, tcpPort)

		result, err = httpClient.Get(appUrl)
		Expect(err).To(HaveOccurred())

		tcpLabel := fmt.Sprintf("tcp-route-sync=%s", tcpPort)
		Eventually(runner.RunKubectlCommand("label", "services", "nginx-ingress-controller", tcpLabel), "10s").Should(gexec.Exit(0))
		Eventually(func() error {
			_, err := httpClient.Get(appUrl)
			return err
		}, "120s", "5s").ShouldNot(HaveOccurred())
		hasPassed = true
	})

})
