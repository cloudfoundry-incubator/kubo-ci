package cloudfoundry_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"integration-tests/test_helpers"
)

var _ = Describe("Deploy workload", func() {

	var (
		nginxSpec        = test_helpers.PathFromRoot("specs/nginx.yml")
		runner           *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
		runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", nginxSpec)
		session.Wait("30s")

		runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	})

	It("exposes routes via CF routers", func() {
		By("deploying application")
		Eventually(runner.RunKubectlCommand("create", "-f", nginxSpec), "60s").Should(gexec.Exit(0))

		serviceName := "test-" + test_helpers.GenerateRandomName()
		appUrl := fmt.Sprintf("http://%s.%s", serviceName, appsDomain)

		By("exposing it via HTTP")
		result, err := http.Get(appUrl)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.StatusCode).NotTo(Equal(200))

		httpLabel := fmt.Sprintf("http-route-sync=%s", serviceName)
		Eventually(runner.RunKubectlCommand("label", "services", "nginx", httpLabel), "10s").Should(gexec.Exit(0))

		timeout := time.Duration(5 * time.Second)
		httpClient := http.Client{
			Timeout: timeout,
		}
		Eventually(func() int {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				return -1
			}
			return result.StatusCode
		}, "120s", "5s").Should(Equal(200))

		By("exposing it via TCP")
		appUrl = fmt.Sprintf("http://%s:%d", tcpRouterDNSName, tcpPort)

		result, err = http.Get(appUrl)
		Expect(err).To(HaveOccurred())

		tcpLabel := fmt.Sprintf("tcp-route-sync=%d", tcpPort)
		Eventually(runner.RunKubectlCommand("label", "services", "nginx", tcpLabel), "10s").Should(gexec.Exit(0))
		Eventually(func() error {
			_, err := http.Get(appUrl)
			return err
		}, "120s", "5s").ShouldNot(HaveOccurred())
	})
})
