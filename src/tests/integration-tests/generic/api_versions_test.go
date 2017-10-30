package generic_test

import (
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("API Versions", func() {
	var (
		runner *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
	})

	It("has RBAC enabled", func() {
		lines := runner.GetOutput("api-versions")

		Expect(lines).To(ContainElement(MatchRegexp("^rbac.*")))
	})

})
