package worker_drain

import (
	. "tests/test_helpers"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/bosh-turbulence/turbulence/incident"
	"github.com/bosh-turbulence/turbulence/incident/selector"
	"github.com/bosh-turbulence/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Worker drain scenarios", func() {

	var (
		deployment     director.Deployment
		err            error
		deploymentName string
		kubectl        *KubectlRunner
		iaas           string
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
		DeploySmorgasbord(kubectl, iaas)
	})

	AfterEach(func() {
		DeleteSmorgasbord(kubectl, iaas)
		kubectl.Teardown()
	})

	Specify("Drain doesn't fail with temporary network issues", func() {
		vmInfos := DeploymentVmsOfType(deployment, WorkerVMType, VMRunningState)
		blockedWorkerID := vmInfos[0].ID

		hellRaiser := TurbulenceClient()
		blockOneWorker := incident.Request{
			Selector: selector.Request{
				Deployment: &selector.NameRequest{
					Name: deploymentName,
				},
				Group: &selector.NameRequest{
					Name: WorkerVMType,
				},
				ID: &selector.IDRequest{
					Values: []string{blockedWorkerID},
				},
			},
			Tasks: tasks.OptionsSlice{
				tasks.FirewallOptions{
					Type:    "Firewall",
					Timeout: "3m",
				},
			},
		}

		By("Recreating all workers successfully")
		dir := NewDirector()
		deployment, err := dir.FindDeployment(deploymentName)
		Expect(err).NotTo(HaveOccurred())
		hellRaiser.CreateIncident(blockOneWorker)
		err = deployment.Recreate(director.NewAllOrInstanceGroupOrInstanceSlug("worker", blockedWorkerID), director.RecreateOpts{Canaries: "0", MaxInFlight: "100%"})
		Expect(err).NotTo(HaveOccurred())
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

})
