package cluster_restart_test

import (
	"tests/test_helpers"
	. "tests/test_helpers"

	"os"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cluster upgrade", func() {
	var (
		deploymentName string
		deployment     director.Deployment
		err            error
		kubectl        *KubectlRunner
		iaas           string
		nginxSpec      string
	)

	BeforeEach(func() {
		director := NewDirector()
		iaas = os.Getenv("IAAS")
		deploymentName = os.Getenv("BOSH_DEPLOYMENT")
		deployment, err = director.FindDeployment(deploymentName)
		Expect(err).NotTo(HaveOccurred())

		kubectl = NewKubectlRunner()
		kubectl.Setup()

		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
		nginxSpec = test_helpers.PathFromRoot("specs/nginx-lb.yml")
		if iaas == "vsphere" {
			nginxSpec = test_helpers.PathFromRoot("specs/nginx-specified-nodeport.yml")
		}
		kubectl.StartKubectlCommand("apply", "-f", nginxSpec).Wait("60s")
		DeploySmorgasbord(kubectl, iaas)
	})

	AfterEach(func() {
		kubectl.StartKubectlCommand("delete", "-f", nginxSpec).Wait("60s")
		DeleteSmorgasbord(kubectl, iaas)
		kubectl.Teardown()
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	Specify("doesn't fail when deployment is recreated", func() {
		dir := NewDirector()
		deployment, err := dir.FindDeployment(deploymentName)
		Expect(err).NotTo(HaveOccurred())
		err = deployment.Recreate(director.AllOrInstanceGroupOrInstanceSlug{}, director.RecreateOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
		WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*5)
	})

	Specify("doesn't fail when deployment is restarted", func() {
		dir := NewDirector()
		deployment, err := dir.FindDeployment(deploymentName)
		Expect(err).NotTo(HaveOccurred())
		err = deployment.Restart(director.AllOrInstanceGroupOrInstanceSlug{}, director.RestartOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
		WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*5)
	})
})
