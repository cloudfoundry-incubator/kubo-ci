package persistence_failure_test

import (
	. "turbulence-tests/test_helpers"

	"fmt"

	"strings"

	"errors"

	"net/http"

	"math/rand"

	"io/ioutil"

	"strconv"

	"time"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Worker failure scenarios", func() {

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
		pvcSpec := PathFromRoot("specs/persistent-volume-claim.yml")
		Eventually(kubectl.RunKubectlCommand("delete", "-f", pvcSpec), "60s").Should(gexec.Exit(0))
		storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
		Eventually(kubectl.RunKubectlCommand("delete", "-f", storageClassSpec), "60s").Should(gexec.Exit(0))
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Specify("K8s applications with persistence keeps their data when node is destroyed", func() {

		By("Deploying the persistent application the value is persisted")

		deployGuestBook(kubectl)

		appAddress := getAppAddress(deployment, kubectl)

		testValue := strconv.Itoa(rand.Int())
		println(testValue)

		postToGuestBook(appAddress, testValue)

		Eventually(func() string {
			return getValueFromGuestBook(appAddress)
		}, "120s", "2s").Should(ContainSubstring(testValue))

		By("Un-deploying the application and re-deploying the data is still available from the persisted source")

		undeployGuestBook(kubectl)
		deployGuestBook(kubectl)

		appAddress = getAppAddress(deployment, kubectl)
		Eventually(func() string {
			return getValueFromGuestBook(appAddress)
		}, "120s", "2s").Should(ContainSubstring(testValue))

		externalId := getExternalId(kubectl, iaas)
		By(fmt.Sprintf("Deleting the node/worker (%s) the persisted data is still available to the application", externalId))
		KillVMById(externalId, iaas)

		fmt.Println(time.Now())
		Eventually(func() string {
			return getValueFromGuestBook(appAddress)
		}, "600s", "2s").Should(ContainSubstring(testValue))
		fmt.Println(time.Now())

		By("Deleting the worker a new worker replaces it")
		Eventually(func() bool { return AllBoshWorkersHaveJoinedK8s(deployment, kubectl) }, 600, 20).Should(BeTrue())

	})

})

func getAppAddress(deployment director.Deployment, kubectl *KubectlRunner) string {
	workerIP := GetWorkerIP(deployment)
	nodePort, err := GetNodePort(kubectl)
	Expect(err).ToNot(HaveOccurred())

	return fmt.Sprintf("%s:%s", workerIP, nodePort)
}

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

func undeployGuestBook(kubectl *KubectlRunner) {
	guestBookSpec := PathFromRoot("specs/pv-guestbook.yml")
	Eventually(kubectl.RunKubectlCommand("delete", "-f", guestBookSpec), "120s").Should(gexec.Exit(0))
}

func deployGuestBook(kubectl *KubectlRunner) {

	guestBookSpec := PathFromRoot("specs/pv-guestbook.yml")
	Eventually(kubectl.RunKubectlCommand("apply", "-f", guestBookSpec), "120s").Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("rollout", "status", "deployment/frontend", "-w"), "120s").Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("rollout", "status", "deployment/redis-master", "-w"), "120s").Should(gexec.Exit(0))

}

func postToGuestBook(address string, testValue string) {

	url := fmt.Sprintf("http://%s/guestbook.php?cmd=set&key=messages&value=%s", address, testValue)
	_, err := http.Get(url)
	Expect(err).ToNot(HaveOccurred())

}

func getValueFromGuestBook(address string) string {

	httpClient := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	url := fmt.Sprintf("http://%s/guestbook.php?cmd=get&key=messages", address)
	response, err := httpClient.Get(url)
	if err != nil {
		return fmt.Sprintf("error occured : %s", err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).ToNot(HaveOccurred())
	return string(bodyBytes)

}

func GetWorkerIP(deployment director.Deployment) string {
	vms := DeploymentVmsOfType(deployment, WorkerVmType, VmRunningState)
	return vms[0].IPs[0]
}

func GetNodePort(kubectl *KubectlRunner) (string, error) {
	output := kubectl.GetOutput("describe", "svc/frontend")

	for i := 0; i < len(output); i++ {
		if output[i] == "NodePort:" {
			nodePort := output[i+2]
			return nodePort[:strings.Index(nodePort, "/")], nil
		}
	}

	return "", errors.New("No nodePort found!")
}
