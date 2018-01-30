package generic_test

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("API Versions", func() {
	var (
		kubectl *KubectlRunner
	)

	BeforeEach(func() {
		kubectl = NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
	})

	It("has RBAC enabled", func() {
		lines := kubectl.GetOutput("api-versions")

		Expect(lines).To(ContainElement(MatchRegexp("^rbac.*")))
	})

})
