package gcp_lb_test

import (
	"fmt"
	"net/http"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Deploy workload", func() {

	It("exposes routes via GCP LBs", func() {

		appUrl := fmt.Sprintf("http://%s", workerAddress)

		timeout := time.Duration(5 * time.Second)
		httpClient := http.Client{
			Timeout: timeout,
		}

		_, err := httpClient.Get(appUrl)
		Expect(err).To(HaveOccurred())

		Eventually(runner.RunKubectlCommand("create", "-f", nginxSpec), "60s").Should(gexec.Exit(0))


		Eventually(func() int {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				return -1
			}
			return result.StatusCode
		}, "120s", "5s").Should(Equal(200))


	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", nginxSpec)
		session.Wait("30s")
	})

})
