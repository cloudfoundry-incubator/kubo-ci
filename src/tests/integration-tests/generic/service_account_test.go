package generic_test

import (
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Service Accounts", func() {
	var (
		runner *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
	})

	It("Should show kube-dns running with the kube-dns service account", func() {
		s := runner.RunKubectlCommandInNamespace("kube-system", "get", "deployment/kube-dns",
			"-o", "jsonpath='{.spec.template.spec.serviceAccountName}'")
		Eventually(s, "15s").Should(gexec.Exit(0))
		Expect(string(s.Out.Contents())).To(Equal("'kube-dns'"))
	})

	It("Should show heapster running with the heapster service account", func() {
		s := runner.RunKubectlCommandInNamespace("kube-system", "get", "deployment/heapster",
			"-o", "jsonpath='{.spec.template.spec.serviceAccountName}'")
		Eventually(s, "15s").Should(gexec.Exit(0))
		Expect(string(s.Out.Contents())).To(Equal("'heapster'"))

		s = runner.RunKubectlCommandInNamespace("kube-system", "get", "deployment/monitoring-influxdb",
			"-o", "jsonpath='{.spec.template.spec.serviceAccountName}'")
		Eventually(s, "15s").Should(gexec.Exit(0))
		Expect(string(s.Out.Contents())).To(Equal("'heapster'"))
	})

	It("Should show kubernetes-dashboard running with the kubernetes-dashboard service account", func() {
		s := runner.RunKubectlCommandInNamespace("kube-system", "get", "deployment/kubernetes-dashboard",
			"-o", "jsonpath='{.spec.template.spec.serviceAccountName}'")
		Eventually(s, "15s").Should(gexec.Exit(0))
		Expect(string(s.Out.Contents())).To(Equal("'kubernetes-dashboard'"))
	})
})
