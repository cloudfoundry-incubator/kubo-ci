package workload_test

import (
	"turbulence-tests/test_helpers"

	"github.com/cppforlife/turbulence/incident"

	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker failure test", func() {

	It("brings back the failed worker VM", func() {
		client := test_helpers.TurbulenceClient()

		limit, _ := selector.NewLimitFromString("1")

		req := incident.Request{
			Tasks: tasks.OptionsSlice{tasks.KillOptions{}},
			Selector: selector.Request{
				Deployment: &selector.NameRequest{Name: "ci-service"},
				Group:      &selector.NameRequest{Name: "worker"},
				ID:         &selector.IDRequest{Limit: limit},
			},
		}

		inc := client.CreateIncident(req)
		inc.Wait()

		Expect(inc.HasTaskErrors()).To(BeFalse())

		Eventually(test_helpers.RunningVmList("worker"), 600, 20).Should(HaveLen(3))
	})
})
