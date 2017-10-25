package persistence_failure_test

import (
	. "tests/test_helpers"

	"fmt"

	"math/rand"

	"strconv"

	"time"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Persistence failure scenarios", func() {

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
		kubectl.CreateNamespace()

		Expect(countRunningWorkers()).To(Equal(3))
		Expect(AllBoshWorkersHaveJoinedK8s(deployment, kubectl)).To(BeTrue())

		storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
		Eventually(kubectl.RunKubectlCommand("create", "-f", storageClassSpec), "60s").Should(gexec.Exit(0))
		pvcSpec := PathFromRoot("specs/persistent-volume-claim.yml")
		Eventually(kubectl.RunKubectlCommand("create", "-f", pvcSpec), "60s").Should(gexec.Exit(0))

	})

	AfterEach(func() {
		UndeployGuestBook(kubectl)
		pvcSpec := PathFromRoot("specs/persistent-volume-claim.yml")
		Eventually(kubectl.RunKubectlCommand("delete", "-f", pvcSpec), "60s").Should(gexec.Exit(0))
		storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
		Eventually(kubectl.RunKubectlCommand("delete", "-f", storageClassSpec), "60s").Should(gexec.Exit(0))
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Specify("K8s applications with persistence keeps their data when node is destroyed", func() {

		By("Deploying the persistent application the value is persisted")

		DeployGuestBook(kubectl)

		appAddress := kubectl.GetAppAddress(deployment, "svc/frontend")

		testValue := strconv.Itoa(rand.Int())
		println(testValue)

		PostToGuestBook(appAddress, testValue)

		Eventually(func() string {
			return GetValueFromGuestBook(appAddress)
		}, "120s", "2s").Should(ContainSubstring(testValue))

		By("Un-deploying the application and re-deploying the data is still available from the persisted source")

		UndeployGuestBook(kubectl)
		DeployGuestBook(kubectl)

		appAddress = kubectl.GetAppAddress(deployment, "svc/frontend")
		Eventually(func() string {
			return GetValueFromGuestBook(appAddress)
		}, "120s", "2s").Should(ContainSubstring(testValue))

		externalId := getExternalId(kubectl, iaas)
		By(fmt.Sprintf("Deleting the node/worker (%s) the persisted data is still available to the application", externalId))
		KillVMById(externalId, iaas)

		fmt.Println(time.Now())
		Eventually(func() string {
			return GetValueFromGuestBook(appAddress)
		}, "600s", "2s").Should(ContainSubstring(testValue))
		fmt.Println(time.Now())

		By("Deleting the worker a new worker replaces it")
		Eventually(func() bool { return AllBoshWorkersHaveJoinedK8s(deployment, kubectl) }, 600, 20).Should(BeTrue())

	})

})

func getExternalId(kubectl *KubectlRunner, iaas string) string {

	var externalId string

	nodeName := kubectl.GetOutput("get", "pods", "-l", "app=redis", "-o", "jsonpath={.items[0].spec.nodeName}")

	switch iaas {
	case "gcp":
		externalId = nodeName[0]
		break
	case "aws":
		externalId = kubectl.GetOutput("get", "nodes", nodeName[0], "-o", "jsonpath={.spec.externalID}")[0]
		break
	case "vsphere":
		externalId = kubectl.GetOutput("get", "nodes", nodeName[0], "-o", "jsonpath={.status.addresses[?(@.type==\"InternalIP\")].address}")[0]
	default:
		Fail(fmt.Sprintf("Unsupported IaaS: %s", iaas))
	}
	return externalId

}
