package persistent_volume_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Guestbook storage", func() {

	var (
		kubectl *KubectlRunner
		iaas    string
	)

	BeforeEach(func() {
		kubectl = NewKubectlRunnerWithDefaultConfig()
		kubectl.CreateNamespace()

		kubeclient, err := NewKubeClient()
		Expect(err).NotTo(HaveOccurred())
		nodes, err := kubeclient.CoreV1().Nodes().List(meta_v1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		providerID := nodes.Items[0].Spec.ProviderID
		iaas = strings.Split(providerID, ":")[0]
		fmt.Printf("iaas = %s", iaas)
	})

	AfterEach(func() {
		UndeployGuestBook(kubectl)
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	Context("when the storage class for the pvc is provided", func() {
		var (
			storageClassSpec string
			pvcSpec          string
		)

		BeforeEach(func() {
			storageClassSpec = PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
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

			DeployGuestBook(kubectl)

			appAddress := kubectl.GetAppAddress("svc/frontend")

			testValue := strconv.Itoa(rand.Int())

			PostToGuestBook(appAddress, testValue)

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "120s", "2s").Should(ContainSubstring(testValue))

			By("Un-deploying the application and re-deploying the data is still available from the persisted source")

			UndeployGuestBook(kubectl)
			DeployGuestBook(kubectl)

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
			if iaas != "gce" {
				Skip("Default Storage Class is only supported by gce.")
			}

			pvcSpec = PathFromRoot("specs/default-persistent-volume-claim.yml")
			Eventually(kubectl.RunKubectlCommand("create", "-f", pvcSpec), "60s").Should(gexec.Exit(0))
		})

		AfterEach(func() {
			Eventually(kubectl.RunKubectlCommand("delete", "-f", pvcSpec), "60s").Should(gexec.Exit(0))
		})

		It("should persist with the default storage class", func() {

			By("Deploying the persistent application the value is persisted")

			DeployGuestBook(kubectl)

			appAddress := kubectl.GetAppAddress("svc/frontend")

			testValue := strconv.Itoa(rand.Int())
			println(testValue)

			PostToGuestBook(appAddress, testValue)

			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "120s", "2s").Should(ContainSubstring(testValue))

			By("Un-deploying the application and re-deploying the data is still available from the persisted source")

			UndeployGuestBook(kubectl)
			DeployGuestBook(kubectl)

			appAddress = kubectl.GetAppAddress("svc/frontend")
			Eventually(func() string {
				return GetValueFromGuestBook(appAddress)
			}, "120s", "2s").Should(ContainSubstring(testValue))

		})
	})
})
