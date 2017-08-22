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

	It("exposes routes via LBs", func() {
		deployNginx := runner.RunKubectlCommand("create", "-f", nginxLBSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))
		serviceIP := ""
		Eventually(func() string {
			getServiceIp := runner.RunKubectlCommand("get", "service", "nginx", "-o", "jsonpath='{.status.loadBalancer.ingress[0].ip}'")
			Eventually(getServiceIp, "60s").Should(gexec.Exit(0))
			serviceIP = string(getServiceIp.Out.Contents())
			serviceIP = serviceIP[1 : len(serviceIP)-1]
			return serviceIP
		}, "120s", "5s").Should(Not(Equal("")))

		appUrl := fmt.Sprintf("http://%s", serviceIP)

		timeout := time.Duration(5 * time.Second)
		httpClient := http.Client{
			Timeout: timeout,
		}

		Eventually(func() int {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "Failed to get response from %s: %v", appUrl, err)
				return -1
			}
			if result.StatusCode != 200 {
				fmt.Fprintf(GinkgoWriter, "Failed to get response from %s: StatusCode %v", appUrl, result.StatusCode)
			}
			return result.StatusCode
		}, "120s", "5s").Should(Equal(200))
	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", nginxLBSpec)
		session.Wait("60s")
	})

})
