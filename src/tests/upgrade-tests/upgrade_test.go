package upgrade_tests_test

import (
	"fmt"
	"net/http"
	"os"

	"time"

	"os/exec"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CFCR Upgrade", func() {
	It("keeps the workload available", func() {
		Expect(true).To(BeTrue())

		By("getting the LB address")
		loadbalancerAddress := ""
		Eventually(func() string {
			loadbalancerAddress = k8sRunner.GetLBAddress("nginx", testconfig.Bosh.Iaas)
			return loadbalancerAddress
		}, "120s", "5s").Should(Not(Equal("")))

		By("monitoring availability")
		doneChannel := make(chan bool)
		go func(doneChannel chan bool) {
			totalCount := 0
			successCount := 0
			for {
				select {
				case <-doneChannel:
					return
				default:
					appUrl := fmt.Sprintf("http://%s", "localhost:8080")

					timeout := time.Duration(45 * time.Second)
					httpClient := http.Client{
						Timeout: timeout,
					}

					result, err := httpClient.Get(appUrl)
					totalCount++
					if err != nil {
						fmt.Fprintf(os.Stdout, "Failed to get response from %s: %v\n", appUrl, err)
					} else if result != nil && result.StatusCode != 200 {
						fmt.Fprintf(os.Stdout, "Failed to get response from %s: StatusCode %v\n", appUrl, result.StatusCode)
					} else {
						successCount++
					}
					fmt.Fprintf(os.Stdout, "Successfully curled server %d out of %d times (%.2f)\n", successCount, totalCount, float64(successCount)/float64(totalCount))
					time.Sleep(time.Second)
				}
			}
		}(doneChannel)

		By("triggering a cfcr-release upgrade")
		exec.Command("%s/scripts/deploy_k8s_instance.sh")
		deployK8sScript := test_helpers.PathFromRoot("scripts/deploy-k8s-instance.sh")
		cmd := exec.Command(deployK8sScript)
		err := cmd.Run()
		Expect(err).NotTo(HaveOccurred())
		close(doneChannel)

		By("waiting on completion of the upgrade")

		By("reporting the availability during the upgrade")

	})

	AfterEach(func() {

	})
})
