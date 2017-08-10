package workload_test

import (
	"turbulence-tests/test_helpers"

	"github.com/cppforlife/turbulence/incident"

	"io"

	"bytes"
	"strings"

	"github.com/cloudfoundry/bosh-cli/cmd"
	"github.com/cloudfoundry/bosh-cli/ui"
	"github.com/cloudfoundry/bosh-utils/logger"
	turbulence_client "github.com/cppforlife/turbulence/client"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker failure test", func() {

	It("makes a single worker VM fail", func() {
		turbulenceRunner := test_helpers.NewTurbulenceRunner()

		incidentJson := test_helpers.PathFromRoot("src/turbulence-tests/resources/incidents/kubo_worker_kill_single.json")
		incident, err := turbulenceRunner.ApplyIncident(incidentJson)
		Expect(err).ToNot(HaveOccurred())
		Expect(incident.ID).ToNot(BeEmpty())

		incidents, lstErr := turbulenceRunner.ListIncidents()
		Expect(lstErr).ToNot(HaveOccurred())
		Expect(incidents).To(ContainElement(incident))

		sameInc, getErr := turbulenceRunner.GetIncidentById(incident.ID)

		Expect(getErr).ToNot(HaveOccurred())
		Expect(sameInc).To(Equal(incident))

	})

	FIt("brings back the failed worker VM", func() {
		config := turbulence_client.NewConfigFromEnv()
		clientLogger := logger.NewLogger(logger.LevelNone)
		client := turbulence_client.NewFactory(clientLogger).New(config)

		limit, _ := selector.NewLimit(1, 1, false)

		req := incident.Request{
			Tasks: tasks.OptionsSlice{
				tasks.KillOptions{},
			},

			Selector: selector.Request{

				Group: &selector.NameRequest{
					Name:  "worker",
					Limit: limit,
				},
			},
		}

		inc := client.CreateIncident(req)
		inc.Wait()

		Expect(inc.HasTaskErrors()).To(BeFalse())

		Eventually(func() []string {
			l := logger.NewLogger(logger.LevelNone)

			stdout := bytes.NewBuffer([]byte{})
			output := io.MultiWriter(GinkgoWriter, stdout)

			stderr := bytes.NewBuffer([]byte{})
			errors := io.MultiWriter(GinkgoWriter, stderr)

			boshUI := ui.NewWriterUI(output, errors, l)
			someUI := ui.NewPaddingUI(boshUI)

			confUI := ui.NewWrappingConfUI(someUI, l)
			confUI.EnableNonInteractive()

			cmdFactory := cmd.NewFactory(cmd.NewBasicDeps(confUI, l))
			cmd, err := cmdFactory.New([]string{"vms"})
			if err != nil {
				return []string{err.Error()}
			}
			cmd.Execute()
			vmTable := stdout.String()

			return Filter(strings.Split(vmTable, "\n"), func(line string) bool {
				return strings.Contains(line, "worker") && strings.Contains(line, "running")
			})

		}, 600, 20).Should(HaveLen(3))
	})
})

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
