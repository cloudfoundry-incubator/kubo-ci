package cloudfoundry_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"integration-tests/test_helpers"

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
	)

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

		Eventually(runner.RunKubectlCommand(
			"create", "secret", "tls", "tls-kubernetes", "--cert", certFile.Name(), "--key", keyFile.Name())).Should(gexec.Exit(0))

		Eventually(runner.RunKubectlCommand(
			"create", "secret", "generic", "kubernetes-service", fmt.Sprintf("--from-literal=host=%s", kubernetesServiceHost), fmt.Sprintf("--from-literal=port=%s", kubernetesServicePort))).Should(gexec.Exit(0))

		Eventually(runner.RunKubectlCommand(
			"create", "-f", ingressSpec), "60s").Should(gexec.Exit(0))

		Eventually(runner.RunKubectlCommand(
			"rollout", "status", "-f", ingressSpec, "-w"), "300s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		Eventually(runner.RunKubectlCommand(
			"delete", "-f", ingressSpec), "60s").Should(gexec.Exit())

		Eventually(runner.RunKubectlCommand(
			"delete", "secret", "tls-kubernetes")).Should(gexec.Exit())

		Eventually(runner.RunKubectlCommand(
			"delete", "secret", "kubernetes-service")).Should(gexec.Exit())

		runner.RunKubectlCommand(
			"delete", "namespace", runner.Namespace()).Wait("60s")
	})

	It("Allows routing via Ingress Controller", func() {
		serviceName := test_helpers.GenerateRandomName()
		appUrl := fmt.Sprintf("http://%s.%s", serviceName, appsDomain)

		By("exposing it via HTTP - " + serviceName)
		result, err := http.Get(appUrl)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.StatusCode).To(Equal(404))

		httpLabel := fmt.Sprintf("http-route-sync=%s", serviceName)
		Eventually(runner.RunKubectlCommand("label", "services", "nginx-ingress-controller", httpLabel), "10s").Should(gexec.Exit(0))

		httpClient := http.Client{
			Timeout: time.Duration(5 * time.Second),
		}
		Eventually(func() int {
			result, err := httpClient.Get(appUrl + "/simple-http-server")
			if err != nil {
				fmt.Println(err)
				return -1
			}
			return result.StatusCode
		}, "120s", "5s").Should(Equal(200))

		result, err = http.Get(appUrl)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.StatusCode).To(Equal(404))

		By("exposing it via TCP")
		appUrl = fmt.Sprintf("http://%s:%s", tcpRouterDNSName, tcpPort)

		result, err = http.Get(appUrl)
		Expect(err).To(HaveOccurred())

		tcpLabel := fmt.Sprintf("tcp-route-sync=%s", tcpPort)
		Eventually(runner.RunKubectlCommand("label", "services", "nginx-ingress-controller", tcpLabel), "10s").Should(gexec.Exit(0))
		Eventually(func() error {
			_, err := http.Get(appUrl)
			return err
		}, "120s", "5s").ShouldNot(HaveOccurred())

	})
})
