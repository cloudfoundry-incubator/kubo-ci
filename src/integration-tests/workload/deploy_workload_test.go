package workload_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Deploy workload", func() {

	It("exposes routes via GCP LBs", func() {

		appUrl := fmt.Sprintf("http://%s:%s", workerAddress, nodePort)

		timeout := time.Duration(5 * time.Second)
		httpClient := http.Client{
			Timeout: timeout,
		}

		_, err := httpClient.Get(appUrl)
		Expect(err).To(HaveOccurred())

		deployNginx := runner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))

		Eventually(func() int {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				return -1
			}
			return result.StatusCode
		}, "120s", "5s").Should(Equal(200))

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
		session.Wait("30s")
	})

})
