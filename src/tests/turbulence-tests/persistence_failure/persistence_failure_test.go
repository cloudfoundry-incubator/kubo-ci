package persistence_failure_test

import (
	. "tests/test_helpers"

	"fmt"

	"math/rand"

	"strconv"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = PersistenceFailureDescribe("Persistence failure scenarios", func() {

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

		storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", testconfig.Iaas))
		Eventually(kubectl.RunKubectlCommand("create", "-f", storageClassSpec), "60s").Should(gexec.Exit(0))
		pvcSpec := PathFromRoot("specs/persistent-volume-claim.yml")
		Eventually(kubectl.RunKubectlCommand("create", "-f", pvcSpec), "60s").Should(gexec.Exit(0))

	})

	AfterEach(func() {
		UndeployGuestBook(kubectl, testconfig.TimeoutScale)
		pvcSpec := PathFromRoot("specs/persistent-volume-claim.yml")
		Eventually(kubectl.RunKubectlCommand("delete", "-f", pvcSpec), "60s").Should(gexec.Exit(0))
		storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", testconfig.Iaas))
		Eventually(kubectl.RunKubectlCommand("delete", "-f", storageClassSpec), "60s").Should(gexec.Exit(0))
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Specify("K8s applications with persistence keeps their data when node is destroyed", func() {
		testValue := strconv.Itoa(rand.Int())

		By("Deploying the persistent application", func() {
			DeployGuestBook(kubectl, testconfig.TimeoutScale)
			appAddress := kubectl.GetAppAddress("svc/frontend")

			PostToGuestBook(appAddress, testValue)

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "120s", "2s").Should(ContainSubstring(testValue))
		})

		By("Un-deploying and re-deploying the app", func() {
			UndeployGuestBook(kubectl, testconfig.TimeoutScale)
			DeployGuestBook(kubectl, testconfig.TimeoutScale)
			appAddress := kubectl.GetAppAddress("svc/frontend")

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "120s", "2s").Should(ContainSubstring(testValue))
		})

		By("Deleting the node/worker with persistent volume", func() {
			redisVMIp := VMIpOfRedis(kubectl)
			appAddress := kubectl.GetAppAddress("svc/frontend")
			vmID, err := BoshIdByIp(deployment, redisVMIp)
			Expect(err).NotTo(HaveOccurred())

			hellRaiser := TurbulenceClient(testconfig.Turbulence)
			killRedisVM := incident.Request{
				Selector: selector.Request{
					ID: &selector.IDRequest{
						Values: []string{vmID},
					},
				},
				Tasks: tasks.OptionsSlice{
					tasks.KillOptions{},
				},
			}

			incident := hellRaiser.CreateIncident(killRedisVM)
			incident.Wait()

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "600s", "2s").Should(ContainSubstring(testValue))
		})

		Eventually(func() bool { return AllBoshWorkersHaveJoinedK8s(deployment, kubectl) }, 600, 20).Should(BeTrue())
	})

})

func VMIpOfRedis(kubectl *KubectlRunner) string {
	nodeName := kubectl.GetOutput("get", "pods", "-l", "app=redis", "-o", "jsonpath={.items[0].spec.nodeName}")
	return kubectl.GetOutput("get", "nodes", nodeName[0], "-o", "jsonpath={.status.addresses[?(@.type==\"InternalIP\")].address}")[0]
}
