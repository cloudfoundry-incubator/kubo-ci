package persistent_volume_test

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Azure file", func() {

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
			if iaas != "azure" {
				Skip("Azure File is only supported by azure.")
			}
			storageClassSpec = PathFromRoot("specs/storage-class-azure-file.yml")
			kubectl.RunKubectlCommandWithTimeout("apply", "-f", storageClassSpec)
			pvcSpec = PathFromRoot("specs/persistent-volume-claim-azure-file.yml")
			kubectl.RunKubectlCommandWithTimeout("apply", "-f", pvcSpec)
		})

		AfterEach(func() {
			kubectl.StartKubectlCommand("delete", "-f", PathFromRoot("specs/azure-file-pod.yml")).Wait(kubectl.TimeoutInSeconds)
			kubectl.StartKubectlCommand("delete", "-f", pvcSpec).Wait(kubectl.TimeoutInSeconds)
			kubectl.StartKubectlCommand("delete", "-f", storageClassSpec).Wait(kubectl.TimeoutInSeconds)
		})

		It("should attach an Azure file volume to a pod", func() {
			kubectl.RunKubectlCommandWithTimeout("create", "-f", PathFromRoot("specs/azure-file-pod.yml"))
			WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*2)
		})
	})

})
