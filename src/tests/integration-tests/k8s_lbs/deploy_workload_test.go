package k8s_lbs_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func getIPAddressFromEchoserver(appURL string) string {

	httpClient := http.Client{
		Timeout: time.Duration(45 * time.Second),
	}

	var result *http.Response
	Eventually(func() int {
		var err error
		result, err = httpClient.Get(appURL)
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "Failed to get response from %s: %v\n", appURL, err)
			return -1
		}
		if result != nil && result.StatusCode != 200 {
			fmt.Fprintf(GinkgoWriter, "Failed to get response from %s: StatusCode %v\n", appURL, result.StatusCode)
		}
		return result.StatusCode
	}, "300s", "45s").Should(Equal(200))

	body, err := ioutil.ReadAll(result.Body)
	Expect(err).NotTo(HaveOccurred())
	re := regexp.MustCompile("client_address=(.*)")

	return re.FindAllStringSubmatch(string(body), -1)[0][0]
}

var _ = K8SLBDescribe("Deploy workload", func() {

	var loadbalancerAddress string
	It("exposes routes via LBs", func() {
		deployNginx := runner.RunKubectlCommand("create", "-f", nginxLBSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))
		loadbalancerAddress = ""
		Eventually(func() string {
			loadbalancerAddress = runner.GetLBAddress("nginx", testconfig.Iaas)
			return loadbalancerAddress
		}, "240s", "5s").Should(Not(Equal("")))

		appUrl := fmt.Sprintf("http://%s", loadbalancerAddress)

		timeout := time.Duration(45 * time.Second)
		httpClient := http.Client{
			Timeout: timeout,
		}

		Eventually(func() int {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "Failed to get response from %s: %v\n", appUrl, err)
				return -1
			}
			if result != nil && result.StatusCode != 200 {
				fmt.Fprintf(GinkgoWriter, "Failed to get response from %s: StatusCode %v\n", appUrl, result.StatusCode)
			}
			return result.StatusCode
		}, "300s", "45s").Should(Equal(200))
	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", nginxLBSpec)
		session.Wait("60s")
	})
})

var _ = K8SLBDescribe("When deploying a loadbalancer", func() {
	var loadbalancerAddress string

	Context("with externalTrafficPolicy to local", func() {
		It("shows a different source client IPs", func() {
			if testconfig.Iaas != "gcp" {
				Skip("Test only valid for GCP")
			}

			deployEchoserver := runner.RunKubectlCommand("create", "-f", echoserverLBSpec)
			Eventually(deployEchoserver, "120s").Should(gexec.Exit(0))
			rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/echoserver", "-w")
			Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))

			loadbalancerAddress = ""
			Eventually(func() string {
				loadbalancerAddress = runner.GetLBAddress("echoserver", testconfig.Iaas)
				return loadbalancerAddress
			}, "240s", "5s").Should(Not(Equal("")))

			appURL := fmt.Sprintf("http://%s", loadbalancerAddress)
			ipAddress := getIPAddressFromEchoserver(appURL)

			runner.RunKubectlCommandWithTimeout("patch", "svc/echoserver", "-p", "{\"spec\":{\"externalTrafficPolicy\":\"Local\"}}")

			Eventually(func() string {
				return getIPAddressFromEchoserver(appURL)
			}, "600s", "20s").Should(Not(BeEquivalentTo(ipAddress)))
		})
	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", echoserverLBSpec)
		session.Wait("60s")
	})
})
