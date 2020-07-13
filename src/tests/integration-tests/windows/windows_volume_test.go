package windows_test

import (
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	storageClassSpec = test_helpers.PathFromRoot("specs/windows/storage-class-vsphere-windows.yml")
	statefulSetSpec  = test_helpers.PathFromRoot("specs/windows/stateful-set-windows.yml")
)

var _ = Describe("When deploying to a Windows worker", func() {
	BeforeEach(func() {
		if !hasWindowsWorkers {
			Skip("skipping Windows tests since no Windows nodes were detected")
		}
		kubectl = test_helpers.NewKubectlRunner()
		kubectl.Setup()
		Eventually(kubectl.StartKubectlCommand("create", "-f", storageClassSpec), "60s").Should(gexec.Exit(0))
		Eventually(kubectl.StartKubectlCommand("create", "-f", statefulSetSpec), "60s").Should(gexec.Exit(0))
		Eventually(kubectl.StartKubectlCommand("rollout", "status", "statefulset/windows-pv"), "600s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		Eventually(kubectl.StartKubectlCommand("delete", "-f", statefulSetSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit())
		Eventually(kubectl.StartKubectlCommand("delete", "pvc", "--all"), kubectl.TimeoutInSeconds).Should(gexec.Exit())
		Eventually(kubectl.StartKubectlCommand("delete", "pv", "--all"), kubectl.TimeoutInSeconds).Should(gexec.Exit())
		Eventually(kubectl.StartKubectlCommand("delete", "-f", storageClassSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit())
		kubectl.Teardown()
	})

	Context("when a file is saved to folder in PV", func() {
		Context("when a pod is recreated", func() {
			It("should be able to read the previous created file from PV", func() {
				By("Save file to pv")
				Eventually(kubectl.StartKubectlCommand("exec", "windows-pv-0", "powershell",
					"New-Item -Path 'c:\\var\\run\\' -Name 'testfile1.txt' -ItemType 'file' -Value 'This is a text string.'"),
					"60s").Should(gexec.Exit(0))

				By("recreate pod")
				Eventually(kubectl.StartKubectlCommand("delete", "po", "windows-pv-0"), "60s").Should(gexec.Exit(0))
				Eventually(kubectl.StartKubectlCommand("rollout", "status", "statefulset/windows-pv"), "120s").Should(gexec.Exit(0))

				By("check file")
				Eventually(kubectl.StartKubectlCommand("exec", "windows-pv-0", "powershell",
					"dir c:\\var\\run\\testfile1.txt"), "60s").Should(gexec.Exit(0))
			})
		})
	})
})
