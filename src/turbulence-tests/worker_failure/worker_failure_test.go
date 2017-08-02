package workload_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
  "turbulence-tests/test_helpers"
)

var _ = Describe("Worker failure test", func() {

	It("makes a single worker VM fail", func() {
    turbulenceRunner := test_helpers.NewTurbulenceRunner()

    incidentId, err := turbulenceRunner.ApplyIncident(test_helpers.PathFromRoot("turbulence/incidents/kubo_worker_kill_single.json"))
    Expect(err).ToNot(HaveOccurred())
    Expect(incidentId).To(Equal("foo"))



		// appUrl := fmt.Sprintf("http://%s:%s", workerAddress, nodePort)
    //
		// timeout := time.Duration(5 * time.Second)
		// httpClient := http.Client{
		// 	Timeout: timeout,
		// }
    //
		// _, err := httpClient.Get(appUrl)
		// Expect(err).To(HaveOccurred())
    //
		// deployNginx := runner.RunKubectlCommand("create", "-f", nginxSpec)
		// Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		// rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		// Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))
    //
		// Eventually(func() int {
		// 	result, err := httpClient.Get(appUrl)
		// 	if err != nil {
		// 		return -1
		// 	}
		// 	return result.StatusCode
		// }, "120s", "5s").Should(Equal(200))

	})

	// AfterEach(func() {
	// 	session := runner.RunKubectlCommand("delete", "-f", nginxSpec)
	// 	session.Wait("30s")
	// })

})
