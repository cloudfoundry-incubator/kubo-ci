package generic_test

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Kubernetes DNS", func() {
	It("Should be able to resolve the internal service DNS name", func() {
		kubectl := NewKubectlRunner()

		Eventually(
			kubectl.StartKubectlCommand("run", "dashboard-lookup", "--image=tutum/dnsutils",
				"--", "nslookup", "kubernetes-dashboard.kube-system.svc.cluster.local"),
		).Should(gexec.Exit(0))

		Eventually(func() ([]string, error) {
			return kubectl.GetOutput("get", "pod", "-l", "job-name=dashboard-lookup", "-o", "jsonpath='{.items[0].status.phase}")
		}, "30s").Should(ConsistOf("Succeeded"))
	})
})
