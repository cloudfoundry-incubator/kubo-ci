package k8s_lbs_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func getIPAddressFromEchoserver(appURL string) string {

	httpClient := http.Client{
		Timeout: time.Duration(5 * time.Second),
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
	}, "60s", "15s").Should(Equal(200))

	body, err := ioutil.ReadAll(result.Body)
	Expect(err).NotTo(HaveOccurred())
	re := regexp.MustCompile("client_address=(.*)")

	return re.FindAllStringSubmatch(string(body), -1)[0][0]
}

var _ = Describe("When deploying a loadbalancer", func() {
	var loadbalancerAddress string

	Context("with externalTrafficPolicy to local", func() {
		It("shows a different source client IPs", func() {
			if iaas != "gce" {
				Skip("Test only valid for GCE")
			}

			deployEchoserver := runner.RunKubectlCommand("create", "-f", echoserverLBSpec)
			Eventually(deployEchoserver, "120s").Should(gexec.Exit(0))
			rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/echoserver", "-w")
			Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))

			loadbalancerAddress = ""
			Eventually(func() string {
				loadbalancerAddress = runner.GetLBAddress("echoserver", iaas)
				return loadbalancerAddress
			}, "240s", "60s").Should(Not(Equal("")))

			appURL := fmt.Sprintf("http://%s", loadbalancerAddress)
			ipAddress := getIPAddressFromEchoserver(appURL)
			segments := strings.Split(ipAddress, ".")

			runner.RunKubectlCommandWithTimeout("patch", "svc/echoserver", "-p", "{\"spec\":{\"externalTrafficPolicy\":\"Local\"}}")
			prefix := segments[0] + "." + segments[1] + "."

			Eventually(func() string {
				return getIPAddressFromEchoserver(appURL)
			}, "600s", "60s").Should(Not(HavePrefix(prefix)))
		})
	})

	AfterEach(func() {
		if iaas == "gce" {
			session := runner.RunKubectlCommand("delete", "-f", echoserverLBSpec)
			session.Wait("60s")
		}
	})
})
