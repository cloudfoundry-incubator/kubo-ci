package generic_test

import (
	"integration-tests/test_helpers"

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

	It("Should have a ulimit of 65536", func() {
		podName := test_helpers.GenerateRandomName()
		session := runner.RunKubectlCommand("run", podName, "--image", "ubuntu", "--restart=Never", "-ti", "--rm", "--", "/bin/bash", "-c", "'ulimit -n'")
		Eventually(session, "5s").Should(gexec.Exit(0))
		ulimit := string(session.Out.Contents())
		Expect(ulimit).To(Equal("65536"))
	})

	// It("Should be able to run kubectl commands within pod", func() {
	// 	podName := test_helpers.GenerateRandomName()
	// 	session := runner.RunKubectlCommand("run", podName, "--image", "pcfkubo/alpine:stable", "--restart=Never", "--image-pull-policy=Always", "-ti", "--rm", "--", "kubectl", "get", "services")
	// 	session.Wait(120)
	// 	Expect(session).To(gexec.Exit(0))
	// })

	// It("Should provide access to the dashboard", func() {
	// 	session := runner.RunKubectlCommand("proxy")
	// 	Eventually(session).Should(gbytes.Say("Starting to serve on"))

	// 	timeout := time.Duration(5 * time.Second)
	// 	httpClient := http.Client{
	// 		Timeout: timeout,
	// 	}

	// 	appUrl := "http://127.0.0.1:8001/ui/"

	// 	Eventually(func() int {
	// 		result, err := httpClient.Get(appUrl)
	// 		if err != nil {
	// 			return -1
	// 		}
	// 		return result.StatusCode
	// 	}, "120s", "5s").Should(Equal(200))

	// 	session.Terminate()
	// })

})
