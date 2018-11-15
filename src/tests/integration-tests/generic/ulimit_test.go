package generic

import (
	"strconv"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ulimit", func() {
	var (
		runner *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
		runner.Setup()
	})

	AfterEach(func() {
		runner.Teardown()
	})

	It("Should have a ulimit at least of 65536", func() {
		podName := test_helpers.GenerateRandomUUID()
		output, err := runner.GetOutput("run", podName, "--image", "pcfkubo/ulimit", "--restart=Never", "-i", "--rm")
		Expect(err).NotTo(HaveOccurred())
		Expect(strconv.Atoi(output[0])).To(BeNumerically(">=", 65536))
	})

})
