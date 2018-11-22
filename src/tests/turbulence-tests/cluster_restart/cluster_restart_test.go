package cluster_restart_test

import (
	. "tests/test_helpers"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cluster upgrade", func() {
	var (
		deployment director.Deployment
		kubectl    *KubectlRunner
	)

	BeforeEach(func() {
		var err error
		director := NewDirector(testconfig.Bosh)
		deployment, err = director.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())

		kubectl = NewKubectlRunner()
		kubectl.Setup()

		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
		DeploySmorgasbord(kubectl, testconfig.Iaas)
	})

	AfterEach(func() {
		DeleteSmorgasbord(kubectl, testconfig.Iaas)
		kubectl.Teardown()
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	Specify("doesn't fail when deployment is recreated", func() {
		dir := NewDirector(testconfig.Bosh)
		deployment, err := dir.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())
		err = deployment.Recreate(director.AllOrInstanceGroupOrInstanceSlug{}, director.RecreateOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
		WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*5)
	})

	Specify("doesn't fail when deployment is restarted", func() {
		dir := NewDirector(testconfig.Bosh)
		deployment, err := dir.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())
		err = deployment.Restart(director.AllOrInstanceGroupOrInstanceSlug{}, director.RestartOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
		WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*5)
	})
})
