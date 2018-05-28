package persistent_volume_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "tests/test_helpers"

	"github.com/onsi/gomega/gexec"
)

var _ = PersistentVolumeDescribe("Guestbook storage", func() {

	var (
		deployment director.Deployment
		kubectl    *KubectlRunner
	)

	BeforeEach(func() {
		var err error
		director := NewDirector(testconfig.Bosh)
		deployment, err = director.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())

		kubectl = NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
		kubectl.CreateNamespace()
	})

	AfterEach(func() {
		UndeployGuestBook(kubectl, testconfig.TimeoutScale)
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Context("when the storage class for the pvc is provided", func() {
		var (
			storageClassSpec string
			pvcSpec          string
		)

		BeforeEach(func() {
			storageClassSpec = PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", testconfig.Iaas))
			Eventually(kubectl.RunKubectlCommand("create", "-f", storageClassSpec), "60s").Should(gexec.Exit(0))
			pvcSpec = PathFromRoot("specs/persistent-volume-claim.yml")
			Eventually(kubectl.RunKubectlCommand("create", "-f", pvcSpec), "60s").Should(gexec.Exit(0))
		})

		AfterEach(func() {
			Eventually(kubectl.RunKubectlCommand("delete", "-f", pvcSpec), "60s").Should(gexec.Exit(0))
			Eventually(kubectl.RunKubectlCommand("delete", "-f", storageClassSpec), "60s").Should(gexec.Exit(0))
		})

		It("should persist when application was undeployed", func() {

			By("Deploying the persistent application the value is persisted")

			DeployGuestBook(kubectl, testconfig.TimeoutScale)

			appAddress := kubectl.GetAppAddress("svc/frontend")

			testValue := strconv.Itoa(rand.Int())
			println(testValue)

			PostToGuestBook(appAddress, testValue)

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "120s", "2s").Should(ContainSubstring(testValue))

			By("Un-deploying the application and re-deploying the data is still available from the persisted source")

			UndeployGuestBook(kubectl, testconfig.TimeoutScale)
			DeployGuestBook(kubectl, testconfig.TimeoutScale)

			appAddress = kubectl.GetAppAddress("svc/frontend")
			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "120s", "2s").Should(ContainSubstring(testValue))

		})
	})

	Context("when the storage class for the pvc is not provided", func() {
		var (
			pvcSpec string
		)

		BeforeEach(func() {
			if testconfig.Iaas != "gcp" {
				Skip("Default Storage Class is only supported by gcp.")
			}

			pvcSpec = PathFromRoot("specs/default-persistent-volume-claim.yml")
			Eventually(kubectl.RunKubectlCommand("create", "-f", pvcSpec), "60s").Should(gexec.Exit(0))
		})

		AfterEach(func() {
			Eventually(kubectl.RunKubectlCommand("delete", "-f", pvcSpec), "60s").Should(gexec.Exit(0))
		})

		It("should persist with the default storage class", func() {

			By("Deploying the persistent application the value is persisted")

			DeployGuestBook(kubectl, testconfig.TimeoutScale)

			appAddress := kubectl.GetAppAddress("svc/frontend")

			testValue := strconv.Itoa(rand.Int())
			println(testValue)

			PostToGuestBook(appAddress, testValue)

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "120s", "2s").Should(ContainSubstring(testValue))

			By("Un-deploying the application and re-deploying the data is still available from the persisted source")

			UndeployGuestBook(kubectl, testconfig.TimeoutScale)
			DeployGuestBook(kubectl, testconfig.TimeoutScale)

			appAddress = kubectl.GetAppAddress("svc/frontend")
			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "120s", "2s").Should(ContainSubstring(testValue))

		})
	})
})
