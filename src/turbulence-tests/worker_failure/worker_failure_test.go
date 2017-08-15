package workload_test

import (
	"turbulence-tests/test_helpers"

	"os/exec"

	"fmt"
	"io"

	"regexp"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const (
	workerVmType   = "worker"
	vmRunningState = "running"
)

var _ = Describe("Worker failure scenarios", func() {
	var deployment director.Deployment
	var countRunningWorkers func() int
	var kubectl *test_helpers.KubectlRunner

	allBoshWorkersHaveJoinedK8s := func() bool {
		cids := []string{}
		for _, vm := range test_helpers.DeploymentVmsOfType(deployment, workerVmType, vmRunningState) {
			cids = append(cids, vm.VMID)
		}

		s := kubectl.RunKubectlCommand("get", "nodes").Wait(10)
		for _, cid := range cids {
			if ok, err := regexp.MatchString(cid+"\\s+Ready", string(s.Out.Contents())); err != nil || !ok {
				return false
			}
		}
		return true
	}

	BeforeEach(func() {
		var err error

		director := test_helpers.NewDirector()
		deployment, err = director.FindDeployment("ci-service")
		Expect(err).NotTo(HaveOccurred())
		countRunningWorkers = test_helpers.CountDeploymentVmsOfType(deployment, workerVmType, vmRunningState)

		kubectl = test_helpers.NewKubectlRunner()

		Expect(countRunningWorkers()).To(Equal(3))
		Expect(allBoshWorkersHaveJoinedK8s()).To(BeTrue())
	})

	Specify("The resurrected worker node joins the k8s cluster", func() {
		By("Deleting a Worker VM")
		killVM(test_helpers.DeploymentVmsOfType(deployment, workerVmType, vmRunningState))
		Eventually(countRunningWorkers, 600, 20).Should(Equal(2))

		By("Expecting the Worker VM to be resurrected")
		Eventually(countRunningWorkers, 600, 20).Should(Equal(3))

		By("Verifying the Worker VM has joined the cluster")
		Eventually(allBoshWorkersHaveJoinedK8s, 120, 20).Should(BeTrue())
	})
})

func killVM(vms []director.VMInfo) {
	vm := vms[0]
	cid := vm.VMID
	cmd := exec.Command("gcloud", "-q", "compute", "instances", "delete", cid)
	io.WriteString(GinkgoWriter, fmt.Sprintf("%#v", cmd))
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 300, 20).Should(gexec.Exit())
	Expect(session.ExitCode()).To(Equal(0))
}
