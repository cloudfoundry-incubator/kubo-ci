package cloud_provider_test

import (
	"tests/test_helpers"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var hostpathMountCloudProviderSpec string

var _ = Describe("Cloud Provider", func() {
	BeforeEach(func() {
		hostpathMountCloudProviderSpec = test_helpers.PathFromRoot("specs/hostpath-mount-cloudprovider.yml")
	})

	Context("on vSphere", func() {
		BeforeEach(func() {
			if val, _ := test_helpers.IaaS(); val != "vsphere" {
				Skip("only supported on vSphere")
			}
		})

		AfterEach(func() {
			kubectl.StartKubectlCommand("delete", "-f", hostpathMountCloudProviderSpec)
		})

		Specify("the worker node's cloud-provider.ini is empty", func() {
			By("Deploying a workload that hostpath mounts the file")
			deploySpec := kubectl.StartKubectlCommand("create", "-f", hostpathMountCloudProviderSpec)
			Eventually(deploySpec, "60s").Should(gexec.Exit(0))
			Eventually(func() string {
				return kubectl.GetPodStatus(kubectl.Namespace(), "hostpath-mount-cloudprovider")
			}, "120s").Should(Equal("Running"))

			By("cat-ing the file in the container and expecting the output to be empty")
			remoteCommand := "cat /cloud-provider.ini"
			session := kubectl.StartKubectlCommand("exec", "hostpath-mount-cloudprovider", "--", "/bin/sh", "-c", remoteCommand)
			Eventually(session, "10s").Should(gexec.Exit(0))
			fileContents := string(session.Out.Contents())
			// Check len so test output doesn't leak creds on failure
			Expect(len(strings.TrimSpace(fileContents))).To(Equal(0))
		})
	})
})
