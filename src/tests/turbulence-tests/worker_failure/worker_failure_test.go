package workload_test

import (
	. "tests/test_helpers"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os"
)

var _ = Describe("Worker failure scenarios", func() {

	var (
		deployment          director.Deployment
		deploymentName      string
		countRunningWorkers func() int
		kubectl             *KubectlRunner
		nginxDaemonSetSpec  = PathFromRoot("specs/nginx-daemonset.yml")
	)

	BeforeEach(func() {
		var err error
		director := NewDirector()
		deploymentName = os.Getenv("BOSH_DEPLOYMENT")
		deployment, err = director.FindDeployment(deploymentName)
		Expect(err).NotTo(HaveOccurred())
		countRunningWorkers = CountDeploymentVmsOfType(deployment, WorkerVMType, VMRunningState)

		kubectl = NewKubectlRunner()
		kubectl.Setup()

		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	AfterEach(func() {
		kubectl.StartKubectlCommand("delete", "-f", nginxDaemonSetSpec).Wait(kubectl.TimeoutInSeconds)
		kubectl.Teardown()
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	Specify("K8s applications are scheduled on the resurrected node", func() {
		By("Deleting the Worker VM")
		hellRaiser := TurbulenceClient()
		killOneWorker := incident.Request{
			Selector: selector.Request{
				Deployment: &selector.NameRequest{
					Name: deploymentName,
				},
				Group: &selector.NameRequest{
					Name: WorkerVMType,
				},
				ID: &selector.IDRequest{
					Limit: selector.MustNewLimitFromString("1"),
				},
			},
			Tasks: tasks.OptionsSlice{
				tasks.KillOptions{},
			},
		}
		incident := hellRaiser.CreateIncident(killOneWorker)
		incident.Wait()
		Eventually(countRunningWorkers, 600, 20).Should(Equal(2))

		By("Waiting for K8s node to go away")
		Eventually(GetReadyNodes, "240s", "5s").Should(HaveLen(2))

		By("Verifying that the Worker VM has joined the K8s cluster")
		Eventually(GetReadyNodes, "450s", "5s").Should(HaveLen(3))

		By("Ensuring a new worker VM has joined the bosh deployment")
		var runningWorkerVms []director.VMInfo
		getRunningWorkerVms := func() []director.VMInfo {
			runningWorkerVms = DeploymentVmsOfType(deployment, WorkerVMType, VMRunningState)
			return runningWorkerVms
		}
		Eventually(getRunningWorkerVms, kubectl.TimeoutInSeconds, "1s").Should(HaveLen(3))

		By("Deploying nginx on 3 nodes")
		Eventually(kubectl.StartKubectlCommand("create", "-f", nginxDaemonSetSpec), kubectl.TimeoutInSeconds/2, "5s").Should(gexec.Exit(0))
		Eventually(kubectl.StartKubectlCommand("rollout", "status", "daemonset/nginx", "-w"), kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))

		By("Verifying nginx got deployed on new node")
		var nodeNames []string
		Eventually(func() []string {
			nodeNames = GetNodeNamesForRunningPods(kubectl)
			return nodeNames
		}, kubectl.TimeoutInSeconds).Should(HaveLen(3))

		_, err := GetNewVmId(runningWorkerVms, nodeNames)
		Expect(err).ToNot(HaveOccurred())
	})
})
