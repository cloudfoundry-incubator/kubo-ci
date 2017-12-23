package master_failure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"tests/config"
	. "tests/test_helpers"

	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
)

var _ = Describe("A single master and etcd failure", func() {

	var testconfig *config.Config

	BeforeSuite(func() {
		var err error
		testconfig, err = config.InitConfig()
		Expect(err).NotTo(HaveOccurred())
	})

	Specify("The cluster is healthy after master is resurrected", func() {
		director := NewDirector(testconfig.Bosh)
		deployment, err := director.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())
		countRunningApiServerOnMaster := CountProcessesOnVmsOfType(deployment, MasterVmType, "kube-apiserver", VmRunningState)

		Expect(countRunningApiServerOnMaster()).To(Equal(1))

		By("Deleting the Master VM")
		hellRaiser := TurbulenceClient(testconfig.Turbulence)
		killOneMaster := incident.Request{
			Selector: selector.Request{
				Deployment: &selector.NameRequest{
					Name: testconfig.Bosh.Deployment,
				},
				Group: &selector.NameRequest{
					Name: MasterVmType,
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
		kubectl := NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)

		By("Waiting for resurrection")
		Eventually(countRunningApiServerOnMaster, "10m", "20s").Should(Equal(1))
		ExpectAllComponentsToBeHealthy(kubectl)

		By("Checking that all nodes are available")
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})
})
