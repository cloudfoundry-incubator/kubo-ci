package k8s_lbs_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = WorkloadDescribe("Deploy workload", func() {

	var loadbalancerAddress string
	It("exposes routes via LBs", func() {
		deployNginx := runner.RunKubectlCommand("create", "-f", nginxLBSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))
		loadbalancerAddress = ""
		Eventually(func() string {
			loadbalancerAddress = runner.GetLBAddress("nginx", testconfig.Bosh.Iaas)
			return loadbalancerAddress
		}, "120s", "5s").Should(Not(Equal("")))

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
		runner.CleanupServiceWithLB(loadbalancerAddress, nginxLBSpec, testconfig.Bosh.Iaas)
	})

})
