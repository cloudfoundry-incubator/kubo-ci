package workload_test

import (
	. "turbulence-tests/test_helpers"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Worker failure scenarios", func() {
	var deployment director.Deployment
	var countRunningWorkers func() int
	var kubectl *KubectlRunner

	BeforeEach(func() {
		var err error

		director := NewDirector()
		deployment, err = director.FindDeployment("ci-service")
		Expect(err).NotTo(HaveOccurred())
		countRunningWorkers = CountDeploymentVmsOfType(deployment, WorkerVmType, VmRunningState)

		kubectl = NewKubectlRunner()

		Expect(countRunningWorkers()).To(Equal(3))
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	Specify("K8s applications are scheduled on the resurrected node", func() {
		By("Deleting the Worker VM")
		vms := DeploymentVmsOfType(deployment, WorkerVmType, VmRunningState)
		KillVM(vms, iaas)
		Eventually(countRunningWorkers, 600, 20).Should(Equal(2))

		By("Verifying that the Worker VM has joined the K8s cluster")
		Eventually(func() bool { return AllBoshWorkersHaveJoinedK8s(deployment, kubectl) }, 600, 20).Should(BeTrue())

		By("Deploying nginx on 3 nodes")
		kubectl.CreateNamespace()
		nginxSpec := PathFromRoot("specs/nginx-specified-nodeport.yml")
		Eventually(kubectl.RunKubectlCommand("create", "-f", nginxSpec)).Should(gexec.Exit(0))
		Eventually(kubectl.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w"), "120s").Should(gexec.Exit(0))

		By("Verifying nginx got deployed on new node")
		nodeNames := GetNodeNamesForRunningPods(kubectl)
		_, err := NewVmId(vms, nodeNames)
		Expect(err).ToNot(HaveOccurred())
	})
	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})
})
