package generic

import (
	"fmt"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var kubectl *test_helpers.KubectlRunner

var _ = Describe("When deploying a pod with service", func() {
	BeforeEach(func() {
		kubectl = test_helpers.NewKubectlRunner()
		kubectl.Setup()
	})
	AfterEach(func() {
		kubectl.Teardown()
	})

	Context("of type HostPort", func() {
		var (
			nginxHostPortSpec = test_helpers.PathFromRoot("specs/nginx-hostport.yml")
		)

		BeforeEach(func() {
			deployNginx := kubectl.StartKubectlCommand("create", "-f", nginxHostPortSpec)
			Eventually(deployNginx, kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
			rolloutWatch := kubectl.StartKubectlCommand("rollout", "status", "deployment/nginx-hostport", "-w")
			Eventually(rolloutWatch, kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
		})

		AfterEach(func() {
			kubectl.StartKubectlCommand("delete", "-f", nginxHostPortSpec).Wait(kubectl.TimeoutInSeconds)
		})

		It("should be able to connect to <node>:<port>", func() {
			hostIP, err := kubectl.GetOutput("get", "pod", "-l", "app=nginx-hostport",
				"-o", "jsonpath='{@.items[0].status.hostIP}'")
			Expect(err).NotTo(HaveOccurred())
			url := fmt.Sprintf("http://%s:40801", hostIP[0])
			session := kubectl.StartKubectlCommand("run", "curl-hostport",
				"--image=gcr.io/cf-pks-golf/tutum/curl", "--restart=Never", "--", "curl", url)
			Eventually(session, "10s").Should(gexec.Exit(0))

			Eventually(func() ([]string, error) {
				return kubectl.GetOutput("get", "pod/curl-hostport", "-o", "jsonpath='{.status.phase}'")
			}, "30s").Should(ConsistOf("Succeeded"))
		})
	})
})
