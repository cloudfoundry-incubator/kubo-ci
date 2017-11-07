package cloudfoundry_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Testing Ingress Controller", func() {

	var (
		config    IngressTestConfig
		runner    *test_helpers.KubectlRunner
		hasPassed bool
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
		runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
		config = InitializeTestConfig(runner)
		
		certFile, _ := ioutil.TempFile(os.TempDir(), "cert")
		_, err := certFile.WriteString(config.tlsKubernetesCert)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(certFile.Name())

		keyFile, _ := ioutil.TempFile(os.TempDir(), "key")
		_, err = keyFile.WriteString(config.tlsKubernetesPrivateKey)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(keyFile.Name())

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
				fmt.Sprintf("--from-literal=host=%s", config.kubernetesServiceHost),
				fmt.Sprintf("--from-literal=port=%s", config.kubernetesServicePort),
			),
			"60s",
		).Should(gexec.Exit(0))

		config.createIngressController()

		Eventually(runner.RunKubectlCommand("rollout", "status", "deployments/default-http-backend", "-w"), "300s").Should(gexec.Exit(0))
		Eventually(runner.RunKubectlCommand("rollout", "status", "deployments/nginx-ingress-controller", "-w"), "300s").Should(gexec.Exit(0))
		Eventually(runner.RunKubectlCommand("rollout", "status", "deployments/simple-http-server", "-w"), "300s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		if hasPassed {
			config.deleteIngressController()

			Eventually(runner.RunKubectlCommand("delete", "-f", config.ingressSpec), "60s").Should(gexec.Exit())
			Eventually(runner.RunKubectlCommand("delete", "secret", "tls-kubernetes")).Should(gexec.Exit())
			Eventually(runner.RunKubectlCommand("delete", "secret", "kubernetes-service")).Should(gexec.Exit())
			Eventually(runner.RunKubectlCommand("delete", "namespace", runner.Namespace()), "60s").Should(gexec.Exit())
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
		appUrl = fmt.Sprintf("http://%s:%s", tcpRouterDNSName, config.tcpPort)

		result, err = httpClient.Get(appUrl)
		Expect(err).To(HaveOccurred())

		tcpLabel := fmt.Sprintf("tcp-route-sync=%s", config.tcpPort)
		Eventually(runner.RunKubectlCommand("label", "services", "nginx-ingress-controller", tcpLabel), "10s").Should(gexec.Exit(0))
		Eventually(func() error {
			_, err := httpClient.Get(appUrl)
			return err
		}, "120s", "5s").ShouldNot(HaveOccurred())
		hasPassed = true
	})

})
