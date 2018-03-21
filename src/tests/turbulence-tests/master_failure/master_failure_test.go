package master_failure_test

import (
	. "tests/test_helpers"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = MasterFailureDescribe("A single master and etcd failure", func() {

	var (
		deployment                    director.Deployment
		kubectl                       *KubectlRunner
		nginxSpec                     = PathFromRoot("specs/nginx.yml")
		countRunningApiServerOnMaster func() int
	)

	BeforeEach(func() {
		var err error
		director := NewDirector(testconfig.Bosh)
		deployment, err = director.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())
		countRunningApiServerOnMaster = CountProcessesOnVmsOfType(deployment, MasterVmType, "kube-apiserver", VmRunningState)

		Expect(countRunningApiServerOnMaster()).To(Equal(1))

		kubectl = NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
		kubectl.CreateNamespace()
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "-f", nginxSpec, "--force", "--grace-period=0")
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace(), "--force", "--grace-period=0")
	})

	Specify("The cluster is healthy after master is resurrected", func() {
		By("Deploying a workload on the k8s cluster")
		Eventually(kubectl.RunKubectlCommand("create", "-f", nginxSpec), "30s", "5s").Should(gexec.Exit(0))
		Eventually(kubectl.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w"), "120s").Should(gexec.Exit(0))

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

		By("Verifying the master VM has restarted")
		var startingMasterVm []director.VMInfo
		getStartingMasterVm := func() []director.VMInfo {
			startingMasterVm = DeploymentVmsOfType(deployment, MasterVmType, VmStartingState)
			return startingMasterVm
		}
		Eventually(getStartingMasterVm, 600, 20).Should(HaveLen(1))

		By("Waiting for resurrection")
		Eventually(countRunningApiServerOnMaster, "10m", "20s").Should(Equal(1))
		ExpectAllComponentsToBeHealthy(kubectl)

		By("Checking that all nodes are available")
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())

		By("Checking for the workload on the k8s cluster")
		session := kubectl.RunKubectlCommand("get", "deployment", "nginx")
		Eventually(session, "120s").Should(gexec.Exit(0))
	})
})
