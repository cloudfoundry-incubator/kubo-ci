package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Main", func() {
	It("should print all IPs used in a director.yml", func() {
		ipExec, err := gexec.Build("vsphere-cleaner/ipcalc/cmd/ipcalc")
		Expect(err).ToNot(HaveOccurred())

		command := exec.Command(ipExec, "fixtures/director.yml")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gbytes.Say("192.36.56.1"))
		Eventually(session).Should(gbytes.Say("192.36.56.2"))
		Eventually(session).ShouldNot(gbytes.Say("192.36.56.3"))
		Eventually(session).ShouldNot(gbytes.Say("192.36.56.4"))
		Eventually(session).ShouldNot(gbytes.Say("192.36.56.5"))
		Eventually(session).Should(gbytes.Say("192.36.56.6"))
		Eventually(session).Should(gbytes.Say("192.36.56.7"))
		Eventually(session).Should(gexec.Exit(0))
	})
})
