package cloudfoundry_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Testing Ingress Controller", func() {

	var (
		ingressConfig IngressTestConfig
		kubectl       *KubectlRunner
	)

	BeforeEach(func() {
		kubectl = NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
		kubectl.RunKubectlCommand("create", "namespace", kubectl.Namespace()).Wait("60s")
		ingressConfig = InitializeIngressTestConfig(kubectl, testconfig.Kubernetes)

		certFile, _ := ioutil.TempFile(os.TempDir(), "cert")
		_, err := certFile.WriteString(ingressConfig.tlsKubernetesCert)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(certFile.Name())

		keyFile, _ := ioutil.TempFile(os.TempDir(), "key")
		_, err = keyFile.WriteString(ingressConfig.tlsKubernetesPrivateKey)
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(keyFile.Name())

		Eventually(
			kubectl.RunKubectlCommand(
				"create", "secret", "tls", "tls-kubernetes",
				"--cert", certFile.Name(),
				"--key", keyFile.Name(),
			),
			"60s",
		).Should(gexec.Exit(0))

		Eventually(
			kubectl.RunKubectlCommand(
				"create", "secret", "generic", "kubernetes-service",
				fmt.Sprintf("--from-literal=host=%s", ingressConfig.kubernetesServiceHost),
				fmt.Sprintf("--from-literal=port=%s", ingressConfig.kubernetesServicePort),
			),
			"60s",
		).Should(gexec.Exit(0))

		ingressConfig.createIngressController()

		Eventually(kubectl.RunKubectlCommand("rollout", "status", "deployments/default-http-backend", "-w"), "300s").Should(gexec.Exit(0))
		Eventually(kubectl.RunKubectlCommand("rollout", "status", "deployments/nginx-ingress-controller", "-w"), "300s").Should(gexec.Exit(0))
		Eventually(kubectl.RunKubectlCommand("rollout", "status", "deployments/simple-http-server", "-w"), "300s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		ingressConfig.deleteIngressController()

		Eventually(kubectl.RunKubectlCommand("delete", "-f", ingressConfig.ingressSpec), "60s").Should(gexec.Exit())
		Eventually(kubectl.RunKubectlCommand("delete", "secret", "tls-kubernetes"), "60s").Should(gexec.Exit())
		Eventually(kubectl.RunKubectlCommand("delete", "secret", "kubernetes-service"), "60s").Should(gexec.Exit())
		Eventually(kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace()), "60s").Should(gexec.Exit())
	})

	It("Allows routing via Ingress Controller", func() {
		serviceName := GenerateRandomName()
		appUrl := fmt.Sprintf("http://%s.%s", serviceName, testconfig.Cf.AppsDomain)
		httpClient := http.Client{
			Timeout: time.Duration(5 * time.Second),
		}

		By("exposing it via HTTP - " + serviceName)
		result, err := httpClient.Get(appUrl)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.StatusCode).To(Equal(404))

		httpLabel := fmt.Sprintf("http-route-sync=%s", serviceName)
		Eventually(kubectl.RunKubectlCommand("label", "services", "nginx-ingress-controller", httpLabel), "10s").Should(gexec.Exit(0))

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
		appUrl = fmt.Sprintf("http://%s:%s", testconfig.Kubernetes.MasterHost, ingressConfig.tcpPort)

		result, err = httpClient.Get(appUrl)
		Expect(err).To(HaveOccurred())

		tcpLabel := fmt.Sprintf("tcp-route-sync=%s", ingressConfig.tcpPort)
		Eventually(kubectl.RunKubectlCommand("label", "services", "nginx-ingress-controller", tcpLabel), "10s").Should(gexec.Exit(0))
		Eventually(func() error {
			_, err := httpClient.Get(appUrl)
			return err
		}, "120s", "5s").ShouldNot(HaveOccurred())
	})

})
