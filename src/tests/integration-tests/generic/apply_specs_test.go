package generic_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"tests/test_helpers"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("check apply-specs errand has run correctly", func() {
	var runner *test_helpers.KubectlRunner

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
	})

	DescribeTable("all deployments have rolled out", func(deploymentName string) {
		session := runner.RunKubectlCommandInNamespace("default", "rollout", "status", fmt.Sprintf("deployment/%s", deploymentName), "-w")
		Eventually(session, "120s").Should(gexec.Exit(0))
		Eventually(session, "120s").Should(gbytes.Say(fmt.Sprintf("deployment \"%s\" successfully rolled out", deploymentName)))
	},
		Entry("frontend", "frontend"),
		Entry("redis-slave", "redis-slave"),
		Entry("redis-master", "redis-master"),
	)
})
