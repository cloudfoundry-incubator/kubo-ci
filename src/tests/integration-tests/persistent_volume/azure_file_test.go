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
			kubectl.RunKubectlCommand("delete", "-f", PathFromRoot("specs/azure-file-pod.yml")).Wait("60s")
			kubectl.RunKubectlCommand("delete", "-f", pvcSpec).Wait("60s")
			kubectl.RunKubectlCommand("delete", "-f", storageClassSpec).Wait("60s")
		})

		It("should attach an Azure file volume to a pod", func() {
			kubectl.RunKubectlCommandWithTimeout("create", "-f", PathFromRoot("specs/azure-file-pod.yml"))
			WaitForPodsToRun(kubectl, "120s")
		})
	})

})
