package k8s_lbs_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Deploy workload", func() {

	var loadbalancerAddress string
	It("exposes routes via LBs", func() {
		deployNginx := kubectl.StartKubectlCommand("create", "-f", nginxLBSpec)
		Eventually(deployNginx, kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		rolloutWatch := kubectl.StartKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, kubectl.TimeoutInSeconds*3).Should(gexec.Exit(0))
		loadbalancerAddress = ""
		Eventually(func() string {
			loadbalancerAddress = kubectl.GetLBAddress("nginx", iaas)
			return loadbalancerAddress
		}, 10*kubectl.TimeoutInSeconds, "5s").Should(Not(Equal("")))

		appUrl := fmt.Sprintf("http://%s", loadbalancerAddress)

		timeout := time.Duration(45 * time.Second)
		httpClient := http.Client{
			Timeout: timeout,
		}

		Eventually(func() int {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "Failed to get response from %s: %v\n", appUrl, err)
				return -1
			}
			if result != nil && result.StatusCode != 200 {
				fmt.Fprintf(GinkgoWriter, "Failed to get response from %s: StatusCode %v\n", appUrl, result.StatusCode)
			}
			return result.StatusCode
		}, "300s", "45s").Should(Equal(200))
	})

	AfterEach(func() {
		kubectl.StartKubectlCommand("delete", "-f", nginxLBSpec).Wait(kubectl.TimeoutInSeconds)
	})
})
