package worker_drain

import (
	. "tests/test_helpers"

	director "github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = WorkerDrainDescribe("Worker drain scenarios", func() {

	var (
		deployment          director.Deployment
		countRunningWorkers func() int
		kubectl             *KubectlRunner
	)

	BeforeEach(func() {
		var err error
		director := NewDirector(testconfig.Bosh)
		deployment, err = director.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())
		countRunningWorkers = CountDeploymentVmsOfType(deployment, WorkerVmType, VmRunningState)

		kubectl = NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
		kubectl.CreateNamespace()

		Expect(countRunningWorkers()).To(Equal(3))
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	AfterEach(func() {
		Expect(countRunningWorkers()).To(Equal(3))
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())

		DeleteSmorgasbord(kubectl, testconfig.Iaas)
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Specify("Workers are able to drain", func() {
		By("Deploying all of the drain types")
		DeploySmorgasbord(kubectl, testconfig.Iaas)

		By("Recreating all workers successfully")
		dir := NewDirector(testconfig.Bosh)
		deployment, err := dir.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())
		err = deployment.Recreate(director.NewAllOrInstanceGroupOrInstanceSlug("worker", ""), director.RecreateOpts{Canaries: "0", MaxInFlight: "100%"})
		Expect(err).NotTo(HaveOccurred())
	})

})
