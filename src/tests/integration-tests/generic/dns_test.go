package generic_test

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Kubernetes DNS", func() {
	var (
		busyboxSpec     string
		kubectl *KubectlRunner
	)

	BeforeEach(func() {
		busyboxSpec = PathFromRoot("specs/busybox.yml")
		kubectl = NewKubectlRunner()
	})

	AfterEach(func(){
		kubectl.RunKubectlCommandWithTimeout("delete", "-n", "kube-system", "-f", busyboxSpec)
	})

	It("Should be able to resolve the internal service DNS name", func() {

		Eventually(kubectl.StartKubectlCommandInNamespace("kube-system", "apply", "-f", busyboxSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		Eventually(func() string {
			return kubectl.GetPodStatus("kube-system", "busybox")
		}, kubectl.TimeoutInSeconds).Should(Equal("Running"))

		session := kubectl.StartKubectlCommandInNamespace("kube-system", "exec", "busybox", "--", "nslookup", "metrics-server.kube-system.svc.cluster.local")
		Eventually(session, kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
	})
})
