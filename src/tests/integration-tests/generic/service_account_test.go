package generic_test

import (
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service Accounts", func() {
	var (
		runner *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
	})

	It("Should show kube-dns running with the kube-dns service account", func() {
		serviceAccount := runner.GetServiceAccount("kube-dns", "kube-system")
		Expect(serviceAccount).To(Equal("'kube-dns'"))
	})

	It("Should show heapster running with the heapster service account", func() {
		serviceAccount := runner.GetServiceAccount("heapster", "kube-system")
		Expect(serviceAccount).To(Equal("'heapster'"))

		serviceAccount = runner.GetServiceAccount("monitoring-influxdb", "kube-system")
		Expect(serviceAccount).To(Equal("'heapster'"))
	})

	It("Should show kubernetes-dashboard running with the kubernetes-dashboard service account", func() {
		serviceAccount := runner.GetServiceAccount("kubernetes-dashboard", "kube-system")
		Expect(serviceAccount).To(Equal("'kubernetes-dashboard'"))
	})
})
