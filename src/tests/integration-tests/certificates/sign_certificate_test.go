package certificates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "tests/test_helpers"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Certificate signing requests", func() {

	var (
		csrSpec string
		kubectl *KubectlRunner
	)

	BeforeEach(func() {
		csrSpec = PathFromRoot("specs/csr.yml")
		kubectl = NewKubectlRunnerWithDefaultConfig()
	})

	AfterEach(func() {
		Eventually(kubectl.RunKubectlCommand("delete", "-f", csrSpec), "60s").Should(gexec.Exit(0))
	})

	Context("When a user creates a csr", func() {
		It("should be signed when the an admin approves it", func() {
			Eventually(kubectl.RunKubectlCommand("apply", "-f", csrSpec), "60s").Should(gexec.Exit(0))

			Eventually(kubectl.RunKubectlCommand("certificate", "approve", "test-csr"), "60s").Should(gexec.Exit(0))
			Eventually(kubectl.GetOutput("get", "csr", "test-csr", "-o", "jsonpath='{.status.certificate}'"), "60s").Should(Not(BeEmpty()))
		})
	})
})
