package volume_test

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("a pod emptyDir volume should be mounted under /var/vcap/data/kubelet", func() {
	kubectl := NewKubectlRunner()

	BeforeEach(func() {
		kubectl.Setup()
	})

	AfterEach(func() {
		kubectl.Teardown()
	})

	Context("when an emptyDir volume has been mounted in a container", func() {
		podSpecPath := PathFromRoot("specs/pod-emptydir.yml")

		BeforeEach(func() {
			kubectl.RunKubectlCommandWithTimeout("apply", "-f", podSpecPath)
		})

		AfterEach(func() {
			kubectl.RunKubectlCommandToDeleteResourceWithPathToFile(kubectl.Namespace(), kubectl.TimeoutInSeconds*3, podSpecPath)
		})

		It("should appear on the host under a /var/vcap/data/kubelet subdirectory", func() {
			WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*3)

			output := kubectl.RunKubectlCommandWithRetry(kubectl.Namespace(), kubectl.TimeoutInSeconds*3, "exec", "emptydir-pod", "--", "sh", "-c", "find /var/search -name find_me.txt")
			Expect(output).Should(ContainSubstring("simple-vol/find_me.txt"))
		})
	})
})
