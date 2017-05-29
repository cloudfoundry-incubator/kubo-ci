package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	It("exits 0", func() {
		command := exec.Command(pathToCLI, "lock-file")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
	})

	It("requires lock file", func() {
		command := exec.Command(pathToCLI)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit())
		Expect(session.ExitCode()).NotTo(Equal(0))
	})

	XIt("asks vsphere to delete vms", func() {
		command := exec.Command(pathToCLI, "")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
	})
})
