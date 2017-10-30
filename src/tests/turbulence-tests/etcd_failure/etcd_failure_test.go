package etcd_failure_test

import (
	. "tests/test_helpers"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Etcd failure scenarios", func() {
	var deployment director.Deployment
	var countRunningEtcd func() int
	var kubectl *KubectlRunner

	BeforeEach(func() {
		var err error

		director := NewDirector()
		deployment, err = director.FindDeployment("ci-service")
		Expect(err).NotTo(HaveOccurred())
		countRunningEtcd = CountDeploymentVmsOfType(deployment, EtcdVmType, VmRunningState)

		kubectl = NewKubectlRunner()
		kubectl.CreateNamespace()

		Expect(countRunningEtcd()).To(Equal(3))
		Expect(AllEtcdHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Specify("Etcd nodes rejoin the cluster and contain up-to-date data", func() {
		By("Deleting the Etcd VM")
		vms := DeploymentVmsOfType(deployment, EtcdVmType, VmRunningState)
		KillVM(vms, iaas)
		Eventually(countRunningEtcd, 600, 20).Should(Equal(2))

		By("Verifying that the Etcd VM has joined the K8s cluster")
		Eventually(func() bool { return AllEtcdHaveJoinedK8s(deployment, kubectl) }, 600, 20).Should(BeTrue())
	})

})
