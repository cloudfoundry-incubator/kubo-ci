package workload_test

import (
	"turbulence-tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker failure test", func() {

	It("makes a single worker VM fail", func() {
		turbulenceRunner := test_helpers.NewTurbulenceRunner()

		incident, err := turbulenceRunner.ApplyIncident(test_helpers.PathFromRoot("turbulence/incidents/kubo_worker_kill_single.json"))
		Expect(err).ToNot(HaveOccurred())
		Expect(incident.ID).ToNot(BeEmpty())

		incidents, lstErr := turbulenceRunner.ListIncidents()
		Expect(lstErr).ToNot(HaveOccurred())
		Expect(incidents).To(ContainElement(incident))

		sameInc, getErr := turbulenceRunner.GetIncidentById(incident.ID)

		Expect(getErr).ToNot(HaveOccurred())
		Expect(sameInc).To(Equal(incident))

	})
})
