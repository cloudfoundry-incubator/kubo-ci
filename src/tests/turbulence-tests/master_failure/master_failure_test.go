package master_failure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"tests/test_helpers"
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/uuid"
	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
)

var _ = Describe("A single master failure", func() {
	Specify("Master VM is resurrected within 10 minutes", func() {
		boshDirector := test_helpers.NewDirector()
		deployment, err := boshDirector.FindDeployment(test_helpers.DeploymentName)
		Expect(err).NotTo(HaveOccurred())
		countRunningMasters := test_helpers.CountDeploymentVmsOfType(deployment, test_helpers.MasterVmType, test_helpers.VmRunningState)

		Expect(countRunningMasters()).To(Equal(2))

		By("Deleting the Master VM")

		hellRaiser := test_helpers.TurbulenceClient()
		killOneMaster := incident.Request{
			Selector: selector.Request{
				Deployment: &selector.NameRequest{
					Name: test_helpers.DeploymentName,
				},
				Group: &selector.NameRequest{
					Name: test_helpers.MasterVmType,
				},
				ID: &selector.IDRequest{
					Limit: selector.MustNewLimitFromString("1"),
				},
			},
			Tasks: tasks.OptionsSlice{
				tasks.KillOptions{},
			},
		}

		incident := hellRaiser.CreateIncident(killOneMaster)
		By("Killing VM")
		incident.Wait()
		By("Waiting for Bosh to recognize dead VM")
		Expect(countRunningMasters()).Should(Equal(1))
		By("Waiting for resurrection")
		Eventually(countRunningMasters, 600, 20).Should(Equal(2))

		sshOpts, privateKey, err := director.NewSSHOpts(uuid.NewGenerator())
		Expect(err).ToNot(HaveOccurred())

		slug, err := director.NewAllOrInstanceGroupOrInstanceSlugFromString(test_helpers.MasterVmType)
		Expect(err).ToNot(HaveOccurred())

		By("Setting up SSH")
		sshResult, err := deployment.SetUpSSH(slug, sshOpts)
		Expect(err).ToNot(HaveOccurred())

		//Verify both hosts, because figuring out which one to verify is too complicated
		for _, host := range sshResult.Hosts {
			By(fmt.Sprintf("Running SSH on %s", host.Host))
			Eventually(func() string {
				output, err := test_helpers.RunSSHCommand(host.Host, 22, sshOpts.Username, privateKey, "curl http://127.0.0.1:8080/healthz")
				Expect(err).ToNot(HaveOccurred())
				return output
			}, "30s", "5s").Should(Equal("ok"))
		}
	})
})

