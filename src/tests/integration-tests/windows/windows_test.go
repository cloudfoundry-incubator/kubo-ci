package windows_test

import (
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var runner *test_helpers.KubectlRunner

var _ = Describe("When deploying to a Windows worker", func() {
	BeforeEach(func() {
		if !hasWindowsWorkers {
			Skip("skipping Windows tests since no Windows nodes were detected")
		}
		runner = test_helpers.NewKubectlRunner()
		runner.Setup()
	})
	AfterEach(func() {
		runner.Teardown()
	})

	var (
		webServerSpec = test_helpers.PathFromRoot("specs/windows/webserver.yml")
	)

	BeforeEach(func() {
		deploy := runner.RunKubectlCommand("create", "-f", webServerSpec)
		Eventually(deploy, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("wait", "--for=condition=ready", "pod/windows-webserver")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		runner.RunKubectlCommand("delete", "-f", webServerSpec)
	})
	It("should be able to fetch logs from the pod", func() {
		Eventually(func() ([]string, error) { return runner.GetOutput("logs", "windows-webserver") }, "31s").Should(ConsistOf("Listening", "at", "http://*:80/"))
	})
})
