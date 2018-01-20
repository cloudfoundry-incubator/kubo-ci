package master_failure_test

import (
	"tests/config"
	. "tests/test_helpers"

	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
	"github.com/onsi/gomega/gexec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("A single master and etcd failure", func() {

	var (
		testconfig 	*config.Config
		kubectl 	*KubectlRunner
		nginxSpec 	= PathFromRoot("specs/nginx.yml")
	)

	BeforeSuite(func() {
		var err error
		testconfig, err = config.InitConfig()
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		kubectl = NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
		kubectl.CreateNamespace()
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "-f", nginxSpec)
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Specify("The cluster is healthy after master is resurrected", func() {
		director := NewDirector(testconfig.Bosh)
		deployment, err := director.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())
		countRunningApiServerOnMaster := CountProcessesOnVmsOfType(deployment, MasterVmType, "kube-apiserver", VmRunningState)

		Expect(countRunningApiServerOnMaster()).To(Equal(1))

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
