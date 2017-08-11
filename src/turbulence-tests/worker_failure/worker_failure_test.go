package workload_test

import (
	"turbulence-tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker failure test", func() {

	It("brings back the failed worker VM", func() {
		incident := test_helpers.KillVms("worker", "1")

		Expect(incident.HasTaskErrors()).To(BeFalse())

		Eventually(test_helpers.RunningVmList("worker"), 600, 20).Should(HaveLen(3))
	})

})
