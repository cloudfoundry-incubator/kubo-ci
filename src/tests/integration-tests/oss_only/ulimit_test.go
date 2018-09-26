package oss_only_test

import (
	"strconv"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	It("Should have a ulimit at least of 65536", func() {
		podName := test_helpers.GenerateRandomUUID()
		output, err := runner.GetOutput("run", podName, "--image", "pcfkubo/ulimit", "--restart=Never", "-i", "--rm")
		Expect(err).NotTo(HaveOccurred())
		Expect(strconv.Atoi(output[0])).To(BeNumerically(">=", 65536))
	})

})
