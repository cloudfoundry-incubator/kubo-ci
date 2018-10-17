package cri_failure_test

import (
	"time"

	. "tests/test_helpers"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	deployment      boshdir.Deployment
	numberOfWorkers int
	director        boshdir.Director
	kubectl         *KubectlRunner
)

var _ = CRIFailureDescribe("A dockerd failure", func() {

	BeforeEach(func() {
		var err error
		director = NewDirector(testconfig.Bosh)
		deployment, err = director.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())

		countRunningWorkers := CountDeploymentVmsOfType(deployment, WorkerVMType, VMRunningState)
		numberOfWorkers = countRunningWorkers()

		kubectl = NewKubectlRunner()
		kubectl.CreateNamespace()
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Specify("The containers continued to run after dockerd restart", func() {

		By("Deploying a workload on the k8s cluster")
		remoteCommand := "while true; do sleep 30; done;"
		Eventually(kubectl.RunKubectlCommand("run", "busybox", "--image=busybox", "--", "/bin/sh", "-c", remoteCommand))
		Eventually(func() string {
			return kubectl.GetPodStatusBySelector(kubectl.Namespace(), "run=busybox")
		}, "120s").Should(Equal("Running"))

		By("Getting the workload's node/bosh.id")
		session := kubectl.RunKubectlCommand("get", "pod", "-l", "run=busybox", "-o", "jsonpath={.items[0].spec.nodeName}")
		Eventually(session, "10s").Should(gexec.Exit(0))
		nodeName := string(session.Out.Contents())

		session = kubectl.RunKubectlCommand("get", "nodes", nodeName, "-o", "jsonpath={.metadata.labels['bosh\\.id']}")
		Eventually(session, "10s").Should(gexec.Exit(0))
		boshID := string(session.Out.Contents())

		By("Killing dockerd")
		killDockerd := incident.Request{
			Selector: selector.Request{
				Deployment: &selector.NameRequest{
					Name: testconfig.Bosh.Deployment,
				},
				Group: &selector.NameRequest{
					Name: WorkerVMType,
				},
				ID: &selector.IDRequest{
					Values: []string{boshID},
				},
			},
			Tasks: tasks.OptionsSlice{
				tasks.KillProcessOptions{
					MonitoredProcessName: "docker",
				},
			},
		}

		createTurbulenceIncident(killDockerd, true, "Killing dockerd")
		Eventually(func() []boshdir.VMInfo {
			return getDockerState(director, deployment, "running")
		}, "60s").ShouldNot(HaveLen(numberOfWorkers))

		By("Waiting for dockerd to restart")
		Eventually(func() []boshdir.VMInfo {
			return getDockerState(director, deployment, "running")
		}, "60s").Should(HaveLen(numberOfWorkers))

		By("Giving Docker time to notify Kubernetes of turbulence")
		time.Sleep(30 * time.Second)

		By("Checking that the containers have not restarted")
		Eventually(kubectl.RunKubectlCommand("get", "pod", "-l", "run=busybox", "-o", "jsonpath={.items[0].status.containerStatuses[0].restartCount}"), "30s").Should(gbytes.Say("0"))
	})
})

func createTurbulenceIncident(request incident.Request, waitForIncident bool, msg string) {
	hellRaiser := TurbulenceClient(testconfig.Turbulence)
	incident := hellRaiser.CreateIncident(request)
	if waitForIncident {
		incident.Wait()
	}
	Expect(incident.HasTaskErrors()).To(BeFalse())
}

func getDockerState(director boshdir.Director, deployment boshdir.Deployment, desiredState string) []boshdir.VMInfo {
	return ProcessesOnVmsOfType(deployment, WorkerVMType, "docker", desiredState)
}
