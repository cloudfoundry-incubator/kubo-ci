package external_traffic_policy_test

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

func getSourceIPFromEchoserver(appURL string) (string, error) {

	httpClient := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	result, err := httpClient.Get(appURL)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Failed to get response from %s: %v\n", appURL, err)
		return "", err
	}
	if result != nil && result.StatusCode != 200 {
		return "", fmt.Errorf("Failed to get response from %s: StatusCode %v\n", appURL, result.StatusCode)
	}

	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("client_address=(.*)")

	return re.FindAllStringSubmatch(string(body), -1)[0][0], nil
}

var _ = Describe("When deploying a loadbalancer", func() {
	var loadbalancerAddress string

	Context("with externalTrafficPolicy to local", func() {
		It("shows a different source client IPs", func() {
			if iaas != "gce" && iaas != "azure" {
				Skip("Test only valid for GCE and Azure")
			}

			deployEchoserver := kubectl.RunKubectlCommand("create", "-f", echoserverLBSpec)
			Eventually(deployEchoserver, kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
			rolloutWatch := kubectl.RunKubectlCommand("rollout", "status", "deployment/echoserver", "-w")
			Eventually(rolloutWatch, kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))

			loadbalancerAddress = ""
			Eventually(func() string {
				loadbalancerAddress = kubectl.GetLBAddress("echoserver", iaas)
				return loadbalancerAddress
			}, "240s", kubectl.TimeoutInSeconds).Should(Not(Equal("")))

			appURL := fmt.Sprintf("http://%s", loadbalancerAddress)
			var ipAddress string
			Eventually(func() error {
				var err error
				ipAddress, err = getSourceIPFromEchoserver(appURL)
				return err
			}, "90s", "15s").Should(Succeed())
			segments := strings.Split(ipAddress, ".")

			kubectl.RunKubectlCommandWithTimeout("patch", "svc/echoserver", "-p", "{\"spec\":{\"externalTrafficPolicy\":\"Local\"}}")
			prefix := segments[0] + "." + segments[1] + "."

			loadbalancerAddress = kubectl.GetLBAddress("echoserver", iaas)
			appURL = fmt.Sprintf("http://%s", loadbalancerAddress)
			// reset cache
			kubectl.RunKubectlCommand("delete", "pods", "--all")

			Eventually(func() string {
				newPrefix, err := getSourceIPFromEchoserver(appURL)
				if err != nil {
					GinkgoWriter.Write([]byte(err.Error()))
				}
				return newPrefix
			}, "600s", kubectl.TimeoutInSeconds).Should(And(Not(BeEmpty()), Not(HavePrefix(prefix))))
		})
	})

	AfterEach(func() {
		if iaas == "gce" || iaas == "azure" {
			kubectl.RunKubectlCommand("delete", "-f", echoserverLBSpec).Wait(kubectl.TimeoutInSeconds)
		}
	})
})

var _ = Describe("When using a NodePort service", func() {
	Context("with externalTrafficPolicy to local", func() {
		It("shows a different source client IPs", func() {
			if iaas != "vsphere" && iaas != "openstack" {
				Skip("Test only valid for vSphere and Openstack")
			}
			deployEchoserver := kubectl.RunKubectlCommand("create", "-f", echoserverNodePortSpec)
			Eventually(deployEchoserver, kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
			rolloutWatch := kubectl.RunKubectlCommand("rollout", "status", "daemonset/echoserver", "-w")
			Eventually(rolloutWatch, kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))

			appURL := fmt.Sprintf("http://%s", kubectl.GetAppAddress("svc/echoserver"))
			var sourceIP string
			Eventually(func() error {
				var err error
				sourceIP, err = getSourceIPFromEchoserver(appURL)
				return err
			}, "90s", "15s").Should(Succeed())
			segments := strings.Split(sourceIP, ".")

			kubectl.RunKubectlCommandWithTimeout("patch", "svc/echoserver", "-p", "{\"spec\":{\"externalTrafficPolicy\":\"Local\"}}")
			prefix := segments[0] + "." + segments[1] + "."

			// reset cache
			kubectl.RunKubectlCommand("delete", "pods", "--all")

			Eventually(func() string {
				newSourceIP, err := getSourceIPFromEchoserver(appURL)
				if err != nil {
					GinkgoWriter.Write([]byte(err.Error()))
				}
				return newSourceIP
			}, "600s", kubectl.TimeoutInSeconds).Should(And(Not(BeEmpty()), Not(HavePrefix(prefix))))
		})
	})

	AfterEach(func() {
		if iaas == "vsphere" && iaas != "openstack" {
			kubectl.RunKubectlCommand("delete", "-f", echoserverNodePortSpec).Wait(kubectl.TimeoutInSeconds)
		}
	})
})
