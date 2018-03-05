package upgrade_tests_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"time"

	"os/exec"
	"tests/test_helpers"

	"github.com/cppforlife/go-patch/patch"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	yaml "gopkg.in/yaml.v2"
)

var loadbalancerAddress string

var _ = Describe("Upgrade components", func() {
	nginxSpec := test_helpers.PathFromRoot("specs/nginx-lb.yml")

	BeforeEach(func() {
		deployNginx := k8sRunner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		k8sRunner.CleanupServiceWithLB(loadbalancerAddress, nginxSpec, testconfig.Bosh.Iaas, testconfig.AWS.IngressGroupID)
	})

	It("upgrades BOSH and CFCR Release", func() {
		upgradeAndMonitorAvailability("scripts/install-bosh.sh", "bosh", 0.99)
		upgradeAndMonitorAvailability("scripts/deploy-k8s-instance.sh", "cfcr-release", 0.99)
	})

	It("upgrades stemcell", func() {
		applyUpdateStemcellVersionOps(filepath.Join(testconfig.CFCR.DeploymentPath, "manifests", "cfcr.yml"), testconfig.CFCR.UpgradeToStemcellVersion)
		upgradeAndMonitorAvailability("scripts/deploy-k8s-instance.sh", "stemcell", 0.99)
	})
})

func applyUpdateStemcellVersionOps(manifestPath, stemcellVersion string) {
	manifestContents, err := ioutil.ReadFile(manifestPath)
	Expect(err).NotTo(HaveOccurred())

	var oldManifest interface{}
	err = yaml.Unmarshal(manifestContents, &oldManifest)
	Expect(err).NotTo(HaveOccurred())

	newManifest, err := patch.ReplaceOp{
		Path:  patch.MustNewPointerFromString("/stemcells/0/version"),
		Value: stemcellVersion,
	}.Apply(oldManifest)
	Expect(err).NotTo(HaveOccurred())

	newManifestContents, err := yaml.Marshal(newManifest)
	Expect(err).NotTo(HaveOccurred())

	err = ioutil.WriteFile(manifestPath, newManifestContents, os.ModePerm)
	Expect(err).NotTo(HaveOccurred())
}

func upgradeAndMonitorAvailability(pathToScript string, component string, requestLossThreshold float64) {
	By("Getting the LB address")
	Eventually(func() string {
		loadbalancerAddress = k8sRunner.GetLBAddress("nginx", testconfig.Bosh.Iaas)
		return loadbalancerAddress
	}, "120s", "5s").Should(Not(Equal("")))

	By("Waiting until LB address resolves")
	Eventually(func() ([]string, error) {
		return net.LookupHost(loadbalancerAddress)
	}, "5m", "5s").ShouldNot(HaveLen(0))

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	close(doneChannel)
	Expect(err).NotTo(HaveOccurred())

	By("Reporting the availability during the upgrade")
	Expect(float64(successCount) / float64(totalCount)).To(BeNumerically(">=", requestLossThreshold))
}
