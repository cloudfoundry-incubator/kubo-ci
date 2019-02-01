package persistence_failure_test

import (
	. "tests/test_helpers"

	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/bosh-turbulence/turbulence/incident"
	"github.com/bosh-turbulence/turbulence/incident/selector"
	"github.com/bosh-turbulence/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Persistence failure scenarios", func() {

	var (
		deployment director.Deployment
		err        error
		kubectl    *KubectlRunner
		iaas       string
	)

	BeforeEach(func() {
		iaas = os.Getenv("IAAS")
		director := NewDirector()
		deployment, err = director.FindDeployment(os.Getenv("BOSH_DEPLOYMENT"))
		Expect(err).NotTo(HaveOccurred())

		kubectl = NewKubectlRunner()
		kubectl.Setup()

		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())

		storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
		Eventually(kubectl.StartKubectlCommand("create", "-f", storageClassSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		pvcSpec := PathFromRoot("specs/persistent-volume-claim.yml")
		Eventually(kubectl.StartKubectlCommand("create", "-f", pvcSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))

	})

	AfterEach(func() {
		UndeployGuestBook(kubectl)
		pvcSpec := PathFromRoot("specs/persistent-volume-claim.yml")
		Eventually(kubectl.StartKubectlCommand("delete", "-f", pvcSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
		Eventually(kubectl.StartKubectlCommand("delete", "-f", storageClassSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		kubectl.Teardown()
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	Specify("K8s applications with persistence keeps their data when node is destroyed", func() {
		testValue := strconv.Itoa(rand.Int())

		By("Deploying the persistent application", func() {
			DeployGuestBook(kubectl)
			appAddress := kubectl.GetAppAddress("svc/frontend")

			Eventually(func() error {
				return PostToGuestBook(appAddress, testValue)
			}, kubectl.TimeoutInSeconds*2, "5s").Should(Succeed())

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, kubectl.TimeoutInSeconds*2, "2s").Should(ContainSubstring(testValue))
		})

		By("Un-deploying and re-deploying the app", func() {
			UndeployGuestBook(kubectl)
			DeployGuestBook(kubectl)
			appAddress := kubectl.GetAppAddress("svc/frontend")

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, kubectl.TimeoutInSeconds*2, "2s").Should(ContainSubstring(testValue))
		})

		By("Deleting the node/worker with persistent volume", func() {
			redisVMIp := VMIpOfRedis(kubectl)
			appAddress := kubectl.GetAppAddress("svc/frontend")
			vmID, err := BoshIdByIp(deployment, redisVMIp)
			Expect(err).NotTo(HaveOccurred())

			hellRaiser := TurbulenceClient()
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
	nodeName, err := kubectl.GetOutput("get", "pods", "-l", "app=redis", "-o", "jsonpath={.items[0].spec.nodeName}")
	Expect(err).NotTo(HaveOccurred())
	output, err := kubectl.GetOutput("get", "nodes", nodeName[0], "-o", "jsonpath={.status.addresses[?(@.type==\"InternalIP\")].address}")
	Expect(err).NotTo(HaveOccurred())
	return output[0]
}
