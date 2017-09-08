package master_failure_test

import (
	. "tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/onsi/gomega/gexec"
	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
)

var _ = Describe("A single master failure", func() {
	var deployment director.Deployment
	var countRunningMasters func() int
	var kubectl *KubectlRunner
	var nginxSpec = PathFromRoot("specs/nginx.yml")

	BeforeEach(func() {
		var err error

		director := NewDirector()
		deployment, err = director.FindDeployment("ci-service")
		Expect(err).NotTo(HaveOccurred())
		countRunningMasters = CountDeploymentVmsOfType(deployment, MasterVmType, VmRunningState)

		kubectl = NewKubectlRunner()
		kubectl.CreateNamespace()

		Expect(countRunningMasters()).To(Equal(3))
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "-f", nginxSpec)
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Specify("K8s applications are scheduled on the resurrected node", func() {
		By("Deleting the Worker VM")
		hellRaiser := TurbulenceClient()
		killOneMaster := incident.Request{
			Selector: selector.Request{
				Deployment: &selector.NameRequest{
					Name: "ci-service",
				},
				Group: &selector.NameRequest{
					Name: "worker",
					Limit: selector.MustNewLimitFromString("1"),
				},
			},
			Tasks: tasks.OptionsSlice{
				tasks.KillOptions{},
			},
		}

		hellRaiser.CreateIncident(killOneMaster)
		vms := DeploymentVmsOfType(deployment, WorkerVmType, VmRunningState)
		Eventually(countRunningMasters, 600, 20).Should(Equal(2))

		By("Verifying that the Worker VM has joined the K8s cluster")
		Eventually(func() bool { return AllBoshWorkersHaveJoinedK8s(deployment, kubectl) }, 600, 20).Should(BeTrue())

		By("Deploying nginx on 3 nodes")
		Eventually(kubectl.RunKubectlCommand("create", "-f", nginxSpec)).Should(gexec.Exit(0))
		Eventually(kubectl.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w"), "120s").Should(gexec.Exit(0))

		By("Verifying nginx got deployed on new node")
		nodeNames := GetNodeNamesForRunningPods(kubectl)
		_, err := NewVmId(vms, nodeNames)
		Expect(err).ToNot(HaveOccurred())
	})

})
