package generic_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "tests/test_helpers"
)

var _ = Describe("Kubernetes DNS", func() {
	It("Should be able to resolve the internal service DNS name", func() {
		kubectl := NewKubectlRunner()

		podName := kubectl.GetPodNameBySelector("kube-system", "k8s-app=metrics-server")

		session := kubectl.RunKubectlCommandInNamespace("kube-system", "exec", podName, "--", "nslookup", "kubernetes-dashboard.kube-system.svc.cluster.local")
		Eventually(session, "10s").Should(gexec.Exit(0))
	})
})
