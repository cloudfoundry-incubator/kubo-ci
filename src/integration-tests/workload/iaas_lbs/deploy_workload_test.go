package iaas_lbs_test

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

		deployNginx := runner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))

		nodePort := ""
		Eventually(func() string {
			getNodePort := runner.RunKubectlCommand("get", "service", "nginx", "-o", "jsonpath='{.spec.ports[0].nodePort}'")
			Eventually(getNodePort, "60s").Should(gexec.Exit(0))
			nodePort = string(getNodePort.Out.Contents())
			nodePort = nodePort[1 : len(nodePort)-1]
			return nodePort
		}, "120s", "5s").Should(Not(Equal("")))

		appUrl := fmt.Sprintf("http://%s:%s", workerAddress, nodePort)

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

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", nginxSpec)
		session.Wait("60s")
	})

})
