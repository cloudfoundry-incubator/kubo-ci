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

var loadbalancerAddress, nginxSpec string
var requestLossThreshold float64

var _ = Describe("Upgrade components", func() {
	BeforeEach(func() {
		nginxSpec = test_helpers.PathFromRoot("specs/nginx-lb.yml")
		if testconfig.Iaas == "vsphere" {
			requestLossThreshold = 0.1
			nginxSpec = test_helpers.PathFromRoot("specs/nginx-specified-nodeport.yml")
		} else {
			requestLossThreshold = 0.99
		}

		deployNginx := k8sRunner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))

		test_helpers.DeploySmorgasbord(k8sRunner, testconfig.Iaas)
	})

	AfterEach(func() {
		test_helpers.DeleteSmorgasbord(k8sRunner, testconfig.Iaas)
		k8sRunner.CleanupServiceWithLB(loadbalancerAddress, nginxSpec, testconfig.Iaas, testconfig.AWS)
		k8sRunner.RunKubectlCommand("delete", "namespace", k8sRunner.Namespace())
	})

	It("upgrades BOSH and CFCR Release", func() {
		upgradeAndMonitorAvailability("scripts/install-bosh.sh", "bosh", requestLossThreshold)
		upgradeAndMonitorAvailability("scripts/deploy-k8s-instance.sh", "cfcr-release", requestLossThreshold)
	})

	It("upgrades stemcell", func() {
		applyUpdateStemcellVersionOps(filepath.Join(testconfig.CFCR.DeploymentPath, "manifests", "cfcr.yml"), testconfig.CFCR.UpgradeToStemcellVersion)
		upgradeAndMonitorAvailability("scripts/deploy-k8s-instance.sh", "stemcell", requestLossThreshold)
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
		if testconfig.Iaas == "vsphere" {
			director := test_helpers.NewDirector(testconfig.Bosh)
			deployment, err := director.FindDeployment(testconfig.Bosh.Deployment)
			Expect(err).NotTo(HaveOccurred())

			vms := test_helpers.DeploymentVmsOfType(deployment, "haproxy", test_helpers.VmRunningState)
			Expect(vms).NotTo(BeEmpty(), "haproxy should be present in the deployment")

			loadbalancerAddress = vms[0].IPs[0]
		} else {
			loadbalancerAddress = k8sRunner.GetLBAddress("nginx", testconfig.Iaas)
		}
		return loadbalancerAddress
	}, "120s", "5s").Should(Not(Equal("")))

	By("Waiting until LB address resolves")
	Eventually(func() ([]string, error) {
		return net.LookupHost(loadbalancerAddress)
	}, "5m", "5s").ShouldNot(HaveLen(0))

	By("Monitoring availability")
	appUrl := fmt.Sprintf("http://%s", loadbalancerAddress)
	doneChannel := make(chan bool)
	totalCount := 0
	successCount := 0
	curlNginx := func() (int, error) {
		httpClient := http.Client{
			Timeout: time.Duration(45 * time.Second),
		}
		ret, err := httpClient.Get(appUrl)
		if err != nil {
			return 0, err
		}
		return ret.StatusCode, err
	}
	Eventually(curlNginx, "5m", "5s").Should(Equal(200))

	go func(doneChannel chan bool, f func() (int, error)) {
		fmt.Fprintf(os.Stdout, "\nStart curling endpoint %s", appUrl)
		for {
			select {
			case <-doneChannel:
				fmt.Fprintf(os.Stdout, "\nDone curling endpoint. Successful response received %d out of %d times (%.2f)", successCount, totalCount, float64(successCount)/float64(totalCount))
				return
			default:
				result, err := f()
				totalCount++
				if err != nil {
					fmt.Fprintf(os.Stdout, "\nFailed to get response from %s: %v", appUrl, err)
				}
				if result == 200 {
					successCount++
				} else {
					fmt.Fprintf(os.Stdout, "\nFailed to get 200 StatusCode from %s. Instead received StatusCode %v", appUrl, result)
				}
				time.Sleep(time.Second)
			}
		}
	}(doneChannel, curlNginx)

	By(fmt.Sprintf("Running %s upgrade", component))
	if testconfig.Iaas == "vsphere" {
		os.Setenv("DEPLOYMENT_OPS_FILE", "add-haproxy.yml")
	}
	script := test_helpers.PathFromRoot(pathToScript)
	cmd := exec.Command(script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	close(doneChannel)
	Expect(err).NotTo(HaveOccurred())

	By("Reporting the availability during the upgrade")
	Expect(float64(successCount) / float64(totalCount)).To(BeNumerically(">=", requestLossThreshold))

	By("Checking that all workloads are running once again")
	test_helpers.CheckSmorgasbord(k8sRunner, "10m")
}
