package upgrade_tests_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var loadbalancerAddress, nginxSpec string
var requestLossThreshold, masterRequestLossThreshold float64

var _ = Describe("Upgrade components", func() {
	BeforeEach(func() {
		nginxSpec = test_helpers.PathFromRoot("specs/nginx-lb.yml")
		if testconfig.Iaas == "vsphere" {
			nginxSpec = test_helpers.PathFromRoot("specs/nginx-specified-nodeport.yml")
		}
		requestLossThreshold = 0.99

		masterRequestLossThreshold = 0.99

		deployNginx := k8sRunner.RunKubectlCommand("create", "-f", nginxSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))

		test_helpers.DeploySmorgasbord(k8sRunner, testconfig.Iaas)
	})

	AfterEach(func() {
		test_helpers.DeleteSmorgasbord(k8sRunner, testconfig.Iaas)
		session := k8sRunner.RunKubectlCommand("delete", "-f", nginxSpec)
		session.Wait("60s")
	})

	It("upgrades CFCR Release", func() {
		upgradeAndMonitorAvailability(os.Getenv("BOSH_DEPLOY_COMMAND"), "cfcr-release", requestLossThreshold)
	})

})

func getvSphereLoadBalancer() *exec.Cmd {
	director := test_helpers.NewDirector(testconfig.Bosh)
	deployment, err := director.FindDeployment(testconfig.Bosh.Deployment)
	Expect(err).NotTo(HaveOccurred())
	content := []byte(`global
maxconn 64000
spread-checks 4
defaults
timeout connect 5000ms
timeout client 50000ms
timeout server 50000ms
listen worker-nodes
bind *:30303
mode tcp
balance roundrobin`)
	tmpfile, err := ioutil.TempFile("", "haproxy-config-")
	Expect(err).NotTo(HaveOccurred())
	_, err = tmpfile.Write(content)
	Expect(err).NotTo(HaveOccurred())
	vms := test_helpers.DeploymentVmsOfType(deployment, test_helpers.WorkerVmType, test_helpers.VmRunningState)
	for i, vm := range vms {
		_, err = tmpfile.Write([]byte(fmt.Sprintf("\n  server worker%d %s check port 10250", i, vm.IPs[0])))
		Expect(err).NotTo(HaveOccurred())
	}
	tmpfile.Close()
	cmd := exec.Command("haproxy", "-f", tmpfile.Name())
	err = cmd.Start()
	Expect(err).NotTo(HaveOccurred())

	loadbalancerAddress = "localhost:30303"

	appURL := fmt.Sprintf("http://%s", loadbalancerAddress)
	Eventually(func() (int, error) {
		return curlURL(appURL)
	}, "30s", "5s").Should(Equal(200))

	return cmd
}

func upgradeAndMonitorAvailability(pathToScript string, component string, requestLossThreshold float64) {
	By("Getting the LB address")
	if testconfig.Iaas == "vsphere" {
		session := getvSphereLoadBalancer()
		defer session.Process.Kill()
	} else {
		Eventually(func() string {
			loadbalancerAddress = k8sRunner.GetLBAddress("nginx", testconfig.Iaas)
			return loadbalancerAddress
		}, "120s", "5s").Should(Not(Equal("")))

		By("Waiting until LB address resolves")
		Eventually(func() ([]string, error) {
			return net.LookupHost(loadbalancerAddress)
		}, "5m", "5s").ShouldNot(HaveLen(0))
	}

	By("Monitoring workload availability")
	appURL := fmt.Sprintf("http://%s", loadbalancerAddress)
	doneChannel := make(chan bool)
	totalCount := 0
	successCount := 0
	Eventually(func() (int, error) {
		return curlURL(appURL)
	}, "5m", "5s").Should(Equal(200))

	go func(doneChannel chan bool, f func(string) (int, error)) {
		fmt.Fprintf(os.Stdout, "\nStart curling endpoint %s", appURL)
		for {
			select {
			case <-doneChannel:
				fmt.Fprintf(os.Stdout, "\nDone curling workload endpoint. Successful response received %d out of %d times (%.2f)", successCount, totalCount, float64(successCount)/float64(totalCount))
				return
			default:
				result, err := f(appURL)
				totalCount++
				if err != nil {
					fmt.Fprintf(os.Stdout, "\nFailed to get response from %s: %v", appURL, err)
				}
				if result == 200 {
					successCount++
				} else {
					fmt.Fprintf(os.Stdout, "\nFailed to get 200 StatusCode from %s. Instead received StatusCode %v", appURL, result)
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
	}(doneChannel, curlURL)

	masterTotalCount := 0
	masterSuccessCount := 0
	masterDoneChannel := make(chan bool)
	if testconfig.UpgradeTests.IncludeMultiAZ {
		By("Monitoring master availability")
		masterCheck := func() error {
			defer GinkgoRecover()

			k8sMasterRunner := test_helpers.NewKubectlRunner()
			session := k8sMasterRunner.RunKubectlCommandInNamespaceSilent(k8sRunner.Namespace(), "describe", "pod", "nginx")
			session.Wait("120s")
			if session.ExitCode() == 0 {
				return nil
			}

			errorMessage, err := ioutil.ReadAll(session.Out)
			if err != nil {
				return err
			}
			return fmt.Errorf("Failed to run kubectl: %s", errorMessage)
		}
		Eventually(masterCheck, "5m", "5s").Should(BeNil())

		go func(doneChannel chan bool, f func() error) {
			fmt.Fprintf(os.Stdout, "\nStart kubectl describe pod\n")
			for {
				select {
				case <-doneChannel:
					fmt.Fprintf(os.Stdout, "\nDone checking Kubernetes master endpoint. Successful response received %d out of %d times (%.2f)", masterSuccessCount, masterTotalCount, float64(masterSuccessCount)/float64(masterTotalCount))
					return
				default:
					err := f()
					masterTotalCount++
					if err != nil {
						fmt.Fprintf(os.Stdout, "\nFailed to get response from %s: %v", appURL, err)
					} else {
						masterSuccessCount++
					}
					time.Sleep(500 * time.Millisecond)
				}
			}
		}(masterDoneChannel, masterCheck)
	}

	By(fmt.Sprintf("Running %s upgrade", component))
	cmd := exec.Command(pathToScript)
	cmd.Dir = test_helpers.PathFromRoot("..")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	close(doneChannel)
	close(masterDoneChannel)
	Expect(err).NotTo(HaveOccurred())

	By("Reporting the workload availability during the upgrade")
	Expect(float64(successCount)/float64(totalCount)).To(BeNumerically(">=", requestLossThreshold), "workload was unavaible during the upgrade")

	if testconfig.UpgradeTests.IncludeMultiAZ {
		By("Reporting the master availability during the upgrade")
		Expect(float64(masterSuccessCount)/float64(masterTotalCount)).To(BeNumerically(">=", masterRequestLossThreshold), "Kubernetes API was unavailable during the upgrade")
	}

	By("Checking that all workloads are running once again")
	test_helpers.WaitForPodsToRun(k8sRunner, "10m")
}

func curlURL(appURL string) (int, error) {
	httpClient := http.Client{
		Timeout: time.Duration(45 * time.Second),
	}
	ret, err := httpClient.Get(appURL)
	if err != nil {
		return 0, err
	}
	return ret.StatusCode, err
}
