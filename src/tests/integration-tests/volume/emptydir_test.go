package volume_test

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
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
			Eventually(kubectl.StartKubectlCommand("apply", "-f", podSpecPath), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		})

		AfterEach(func() {
			Eventually(kubectl.StartKubectlCommand("delete", "-f", podSpecPath), kubectl.TimeoutInSeconds).Should(gexec.Exit())
		})

		It("should appear on the host under a /var/vcap/data/kubelet subdirectory", func() {
			WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*3)

			Eventually(kubectl.StartKubectlCommand("exec", "emptydir-pod", "--", "sh", "-c", "[ $(find /var/search -name find_me.txt | wc -l) -eq '1' ]"), kubectl.TimeoutInSeconds*3).Should(gexec.Exit(0))
		})
	})
})
