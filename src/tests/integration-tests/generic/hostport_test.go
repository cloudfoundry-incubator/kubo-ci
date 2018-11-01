package generic

import (
	"fmt"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var runner *test_helpers.KubectlRunner

var _ = Describe("When deploying a pod with service", func() {
	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
		runner.Setup()
	})
	AfterEach(func() {
		runner.Teardown()
	})

	Context("of type HostPort", func() {
		var (
			nginxHostPortSpec = test_helpers.PathFromRoot("specs/nginx-hostport.yml")
		)

		BeforeEach(func() {
			deployNginx := runner.RunKubectlCommand("create", "-f", nginxHostPortSpec)
			Eventually(deployNginx, "60s").Should(gexec.Exit(0))
			rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx-hostport", "-w")
			Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))
		})

		AfterEach(func() {
			runner.RunKubectlCommand("delete", "-f", nginxHostPortSpec)
		})
		It("should be able to connect to <node>:<port>", func() {
			hostIP, err := runner.GetOutput("get", "pod", "-l", "app=nginx-hostport",
				"-o", "jsonpath='{@.items[0].status.hostIP}'")
			Expect(err).NotTo(HaveOccurred())
			url := fmt.Sprintf("http://%s:40801", hostIP)
			session := runner.RunKubectlCommand("run", "curl-hostport",
				"--image=tutum/curl", "--restart=Never", "--", "curl", url)
			Eventually(session, "10s").Should(gexec.Exit(0))
		})
	})
})
