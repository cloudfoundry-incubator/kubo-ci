package workload_test

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var ipAddressError = "No IP address found for service"

func getServiceIP() string {
	timeout := time.After(60 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	validIP := regexp.MustCompile(numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock)

	for {
		select {
		case <-timeout:
			return ipAddressError
		case <-tick:
			getServiceIp := runner.RunKubectlCommand("get", "service", "nginx", "-o", "jsonpath='{.status.loadBalancer.ingress[0].ip}'")
			Eventually(getServiceIp, "5s").Should(gexec.Exit())
			serviceIP := string(getServiceIp.Out.Contents())
			// Remove quotes from the IP address
			serviceIP = serviceIP[1 : len(serviceIP)-1]
			if validIP.MatchString(serviceIP) {
				return serviceIP
			}
		}
	}
}

var _ = Describe("Deploy workload", func() {

	It("exposes routes via GCP LBs", func() {

		deployNginx := runner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))

		serviceIP := getServiceIP()
		Expect(serviceIP).To(Not(Equal(ipAddressError)))
		appUrl := fmt.Sprintf("http://%s", serviceIP)

		timeout := time.Duration(5 * time.Second)
		httpClient := http.Client{
			Timeout: timeout,
		}

		Eventually(func() string {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				return err.Error()
			}
			return result.Status
		}, "120s", "5s").Should(Equal("200 OK"))

	})

	It("allows access to pod logs", func() {

		deployNginx := runner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))

		getPodName := runner.RunKubectlCommand("get", "pods", "-o", "jsonpath='{.items[0].metadata.name}'")
		Eventually(getPodName, "5s").Should(gexec.Exit())
		podName := string(getPodName.Out.Contents())

		getLogs := runner.RunKubectlCommand("logs", podName)
		Eventually(getLogs, "5s").Should(gexec.Exit())
		logContent := string(getLogs.Out.Contents())
		// nginx pods do not log much, unless there is an error we should see an empty string as a result
		Expect(logContent).To(Equal(""))

	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", nginxSpec)
		session.Wait("60s")
	})

})
