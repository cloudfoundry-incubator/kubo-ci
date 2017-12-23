package cloudfoundry_test

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Deploy workload", func() {

	var (
		tcpPort   string
		nginxSpec = PathFromRoot("specs/nginx.yml")
		kubectl   *KubectlRunner
	)

	BeforeEach(func() {
		if testconfig.Kubernetes.MasterPort == 0 {
			Fail("Please ensure k8s master port is defined in the test config")
		}

		tcpPort = strconv.Itoa(testconfig.Kubernetes.MasterPort + 10)

		kubectl = NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
		kubectl.RunKubectlCommand(
			"create", "namespace", kubectl.Namespace()).Wait("60s")

		Eventually(kubectl.RunKubectlCommand(
			"create", "-f", nginxSpec), "60s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		Eventually(kubectl.RunKubectlCommand(
			"delete", "-f", nginxSpec), "60s").Should(gexec.Exit())

		kubectl.RunKubectlCommand(
			"delete", "namespace", kubectl.Namespace()).Wait("60s")
	})

	It("exposes routes via CF routers", func() {
		serviceName := kubectl.Namespace()
		appUrl := fmt.Sprintf("http://%s.%s", serviceName, testconfig.Cf.AppsDomain)
		httpClient := http.Client{
			Timeout: time.Duration(5 * time.Second),
		}

		By("exposing it via HTTP")
		result, err := httpClient.Get(appUrl)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.StatusCode).To(Equal(404))

		httpLabel := fmt.Sprintf("http-route-sync=%s", serviceName)
		Eventually(kubectl.RunKubectlCommand("label", "services", "nginx", httpLabel), "10s").Should(gexec.Exit(0))

		Eventually(func() int {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				fmt.Println(err)
				return -1
			}
			return result.StatusCode
		}, "120s", "5s").Should(Equal(200))

		By("exposing it via TCP")
		appUrl = fmt.Sprintf("http://%s:%s", testconfig.Kubernetes.MasterHost, tcpPort)

		result, err = httpClient.Get(appUrl)
		Expect(err).To(HaveOccurred())

		tcpLabel := fmt.Sprintf("tcp-route-sync=%s", tcpPort)
		Eventually(kubectl.RunKubectlCommand("label", "services", "nginx", tcpLabel), "10s").Should(gexec.Exit(0))
		Eventually(func() error {
			_, err := httpClient.Get(appUrl)
			return err
		}, "120s", "5s").ShouldNot(HaveOccurred())
	})
})
