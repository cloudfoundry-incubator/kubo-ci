package generic_test

import (
	"integration-tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Kubectl", func() {
	var (
		runner *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
		runner.RunKubectlCommand(
			"create", "namespace", runner.Namespace()).Wait("60s")
	})

	AfterEach(func() {
		runner.RunKubectlCommand(
			"delete", "namespace", runner.Namespace()).Wait("60s")
	})

	It("Should be able to run kubectl commands within pod", func() {
		podName := test_helpers.GenerateRandomName()
		session := runner.RunKubectlCommand("run", podName, "--image", "pcfkubo/alpine:stable", "--restart=Never", "--image-pull-policy=Always", "-ti", "--rm", "--", "kubectl", "get", "services")
		session.Wait(120)
		Expect(session).To(gexec.Exit(0))
	})

})
