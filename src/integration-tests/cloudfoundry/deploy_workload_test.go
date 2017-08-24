package cloudfoundry_test

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"integration-tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Deploy workload", func() {

	var (
		tcpPort   string
		nginxSpec = test_helpers.PathFromRoot("specs/nginx.yml")
		runner    *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		tcpPort = os.Getenv("WORKLOAD_TCP_PORT")
		if tcpPort == "" {
			Fail("Correct WORKLOAD_TCP_PORT has to be set")
		}

		runner = test_helpers.NewKubectlRunner()
		runner.RunKubectlCommand(
			"create", "namespace", runner.Namespace()).Wait("60s")

		Eventually(runner.RunKubectlCommand(
			"create", "-f", nginxSpec), "60s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		Eventually(runner.RunKubectlCommand(
			"delete", "-f", nginxSpec), "60s").Should(gexec.Exit())

		runner.RunKubectlCommand(
			"delete", "namespace", runner.Namespace()).Wait("60s")
	})

	It("exposes routes via CF routers", func() {
		serviceName := runner.Namespace()
		appUrl := fmt.Sprintf("http://%s.%s", serviceName, appsDomain)
		httpClient := http.Client{
			Timeout: time.Duration(5 * time.Second),
		}

		By("exposing it via HTTP")
		result, err := httpClient.Get(appUrl)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.StatusCode).To(Equal(404))

		httpLabel := fmt.Sprintf("http-route-sync=%s", serviceName)
		Eventually(runner.RunKubectlCommand("label", "services", "nginx", httpLabel), "10s").Should(gexec.Exit(0))

		Eventually(func() int {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				fmt.Println(err)
				return -1
			}
			return result.StatusCode
		}, "120s", "5s").Should(Equal(200))

		By("exposing it via TCP")
		appUrl = fmt.Sprintf("http://%s:%s", tcpRouterDNSName, tcpPort)

		result, err = httpClient.Get(appUrl)
		Expect(err).To(HaveOccurred())

		tcpLabel := fmt.Sprintf("tcp-route-sync=%s", tcpPort)
		Eventually(runner.RunKubectlCommand("label", "services", "nginx", tcpLabel), "10s").Should(gexec.Exit(0))
		Eventually(func() error {
			_, err := httpClient.Get(appUrl)
			return err
		}, "120s", "5s").ShouldNot(HaveOccurred())
	})
})
