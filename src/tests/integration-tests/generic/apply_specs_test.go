package generic_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"tests/test_helpers"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("apply-specs errand has run", func() {
	var runner *test_helpers.KubectlRunner

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
	})

	It("Should be able to run kubectl commands within pod", func() {

		session := runner.RunKubectlCommandInNamespace("default", "rollout", "status", "deployment/frontend", "-w")
		Eventually(session, "120s").Should(gexec.Exit(0))
		Eventually(session, "120s").Should(gbytes.Say("deployment \"frontend\" successfully rolled out"))
	})
	

})
