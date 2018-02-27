package upgrade_tests_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"time"

	"os/exec"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const (
	// There should only be at most one connection failure for each worker rolling
	CONNECTION_FAILURE_THRESHOLD = 3
)

var loadbalancerAddress string

var _ = Describe("Upgrade components", func() {
	nginxSpec := test_helpers.PathFromRoot("specs/nginx-lb.yml")

	BeforeEach(func() {
		deployNginx := k8sRunner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		k8sRunner.CleanupServiceWithLB(loadbalancerAddress, nginxSpec, testconfig.Bosh.Iaas)
	})

	It("Upgrades CFCR Release", func() {
		upgradeAndMonitorAvailability("scripts/deploy-k8s-instance.sh", "cfcr-release")
	})

	XIt("Upgrades stemcell", func() {
		upgradeAndMonitorAvailability("scripts/upgrade-stemcell.sh", "stemcell")
	})

})

func upgradeAndMonitorAvailability(pathToScript string, component string) {
	By("Getting the LB address")
	Eventually(func() string {
		loadbalancerAddress = k8sRunner.GetLBAddress("nginx", testconfig.Bosh.Iaas)
		return loadbalancerAddress
	}, "120s", "5s").Should(Not(Equal("")))

	By("Monitoring availability")
	doneChannel := make(chan bool)
	totalCount := 0
	successCount := 0
	go func(doneChannel chan bool) {
		for {
			select {
			case <-doneChannel:
				return
			default:
				appUrl := fmt.Sprintf("http://%s", loadbalancerAddress)

				timeout := time.Duration(45 * time.Second)
				httpClient := http.Client{
					Timeout: timeout,
				}

				result, err := httpClient.Get(appUrl)
				totalCount++
				if err != nil {
					fmt.Fprintf(os.Stdout, "\nFailed to get response from %s: %v", appUrl, err)
				} else if result != nil && result.StatusCode != 200 {
					fmt.Fprintf(os.Stdout, "\nFailed to get response from %s: StatusCode %v", appUrl, result.StatusCode)
				} else {
					successCount++
				}
				fmt.Fprintf(os.Stdout, "\nSuccessfully curled server %d out of %d times (%.2f)", successCount, totalCount, float64(successCount)/float64(totalCount))
				time.Sleep(time.Second)
			}
		}
	}(doneChannel)

	By(fmt.Sprintf("Running %s upgrade", component))
	script := test_helpers.PathFromRoot(pathToScript)
	cmd := exec.Command(script)
	pwd := os.Getenv("PWD")
	dirs := strings.Split(pwd, "/")
	cmd.Dir = strings.Join(dirs[:4], "/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	close(doneChannel)
	Expect(err).NotTo(HaveOccurred())

	By("Reporting the availability during the upgrade")
	Expect(totalCount - successCount).To(BeNumerically("<=", CONNECTION_FAILURE_THRESHOLD))
}
