package volume_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Guestbook storage", func() {

	var (
		kubectl *KubectlRunner
		iaas    string
	)

	BeforeEach(func() {
		kubectl = NewKubectlRunner()
		kubectl.Setup()

		var err error
		iaas, err = IaaS()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		kubectl.Teardown()
	})

	Context("when the storage class for the pvc is provided", func() {
		var (
			storageClassSpec string
			pvcSpec          string
		)

		BeforeEach(func() {
			storageClassSpec = PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
			Eventually(kubectl.StartKubectlCommand("apply", "-f", storageClassSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
			pvcSpec = PathFromRoot("specs/persistent-volume-claim.yml")
			Eventually(kubectl.StartKubectlCommand("apply", "-f", pvcSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		})

		AfterEach(func() {
			UndeployGuestBook(kubectl)
			Eventually(kubectl.StartKubectlCommand("delete", "-f", pvcSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
			Eventually(kubectl.StartKubectlCommand("delete", "-f", storageClassSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		})

		It("should persist when application was undeployed", func() {

			By("Deploying the persistent application the value is persisted")

			DeployGuestBook(kubectl)

			appAddress := kubectl.GetAppAddress("svc/frontend")

			testValue := strconv.Itoa(rand.Int())

			Eventually(func() error {
				return PostToGuestBook(appAddress, testValue)
			}, kubectl.TimeoutInSeconds*5, "5s").Should(Succeed())

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, kubectl.TimeoutInSeconds*2, "2s").Should(ContainSubstring(testValue))

			By("Un-deploying the application and re-deploying the data is still available from the persisted source")

			UndeployGuestBook(kubectl)
			time.Sleep(30 * time.Second)
			DeployGuestBook(kubectl)

			By("Getting the value from application")
			appAddress = kubectl.GetAppAddress("svc/frontend")
			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, kubectl.TimeoutInSeconds*2, "2s").Should(ContainSubstring(testValue))

		})
	})

	Context("when the storage class for the pvc is not provided", func() {
		var (
			pvcSpec string
		)

		BeforeEach(func() {
			if iaas != "gce" {
				Skip("Default Storage Class is only supported by gce.")
			}

			pvcSpec = PathFromRoot("specs/default-persistent-volume-claim.yml")
			Eventually(kubectl.StartKubectlCommand("create", "-f", pvcSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		})

		AfterEach(func() {
			if iaas == "gce" {
				UndeployGuestBook(kubectl)
				Eventually(kubectl.StartKubectlCommand("delete", "-f", pvcSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
			}
		})

		It("should persist with the default storage class", func() {

			By("Deploying the persistent application the value is persisted")

			DeployGuestBook(kubectl)

			appAddress := kubectl.GetAppAddress("svc/frontend")

			testValue := strconv.Itoa(rand.Int())
			println(testValue)

			Eventually(func() error {
				return PostToGuestBook(appAddress, testValue)
			}, kubectl.TimeoutInSeconds*2, "5s").Should(Succeed())

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, kubectl.TimeoutInSeconds*2, "2s").Should(ContainSubstring(testValue))

			By("Un-deploying the application and re-deploying the data is still available from the persisted source")

			UndeployGuestBook(kubectl)
			DeployGuestBook(kubectl)

			appAddress = kubectl.GetAppAddress("svc/frontend")
			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, kubectl.TimeoutInSeconds*2, "2s").Should(ContainSubstring(testValue))

		})
	})
})
