package etcd_failure_test

import (
	. "tests/test_helpers"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cppforlife/turbulence/client"
)

var _ = Describe("Etcd failure scenarios", func() {
	var deployment director.Deployment
	var countRunningEtcd func() int
	var kubectl *KubectlRunner
	var director director.Director

	BeforeEach(func() {
		var err error

		director = NewDirector()
		deployment, err = director.FindDeployment(DeploymentName)
		Expect(err).NotTo(HaveOccurred())
		countRunningEtcd = CountDeploymentVmsOfType(deployment, EtcdVmType, VmRunningState)

		kubectl = NewKubectlRunner()
		kubectl.CreateNamespace()

		Expect(countRunningEtcd()).To(Equal(1))
		Expect(AllEtcdNodesAreHealthy(deployment, kubectl)).To(BeTrue())
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Specify("Etcd nodes rejoin the cluster", func() {
		By("Deleting the Etcd VM", func() {
			turbulenceClient := TurbulenceClient()
			killOneEtcd := incident.Request{
				Selector: selector.Request{
					Deployment: &selector.NameRequest{
						Name: DeploymentName,
					},
					Group: &selector.NameRequest{
						Name: EtcdVmType,
					},
					ID: &selector.IDRequest{
						Limit: selector.MustNewLimitFromString("1"),
					},
				},
				Tasks: tasks.OptionsSlice{
					tasks.KillOptions{},
				},
			}
			incident := turbulenceClient.CreateIncident(killOneEtcd)
			incident.Wait()
			Eventually(countRunningEtcd, 600, 20).Should(Equal(0))
		})

		By("Waiting for Bosh Resurrection", func() {
			Eventually(countRunningEtcd, 600, 20).Should(Equal(1))
			Eventually(func() bool { return AllEtcdNodesAreHealthy(deployment, kubectl) }, 600, 20).Should(BeTrue())
		})
	})

})
