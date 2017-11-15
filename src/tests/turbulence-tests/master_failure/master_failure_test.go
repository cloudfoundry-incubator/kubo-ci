package master_failure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"tests/test_helpers"

	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
)

var _ = Describe("A single master and etcd failure", func() {
	Specify("The cluster is healthy after master is resurrected", func() {
		boshDirector := test_helpers.NewDirector()
		deployment, err := boshDirector.FindDeployment(test_helpers.DeploymentName)
		Expect(err).NotTo(HaveOccurred())
		countRunningApiServerOnMaster := test_helpers.CountProcessesOnVmsOfType(deployment, test_helpers.MasterVmType, "kubernetes-api", test_helpers.VmRunningState)

		Expect(countRunningApiServerOnMaster()).To(Equal(1))

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
		incident.Wait()
		Expect(countRunningApiServerOnMaster()).Should(Equal(0))
		kubectl := test_helpers.NewKubectlRunner()

		By("Waiting for resurrection")
		Eventually(countRunningApiServerOnMaster, "10m", "20s").Should(Equal(1))
		test_helpers.ExpectAllComponentsToBeHealthy(kubectl)

		By("Checking that all nodes are available")
		Expect(test_helpers.AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})
})
