package windows_test

import (
	"fmt"
	"math/rand"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	kubectl       *test_helpers.KubectlRunner
	webServerSpec = test_helpers.PathFromRoot("specs/windows/webserver.yml")
)

var _ = Describe("When deploying to a Windows worker", func() {

	It("has functional pod networking", func() {
		setupNS()
		defer teardown()

		deploy := kubectl.StartKubectlCommand("create", "-f", webServerSpec)
		Eventually(deploy, "60s").Should(gexec.Exit(0))
		Eventually(kubectl.StartKubectlCommand("wait", "--timeout=120s",
			"--for=condition=ready", "pod/windows-webserver"), "120s").Should(gexec.Exit(0))

		By("should be able to fetch logs from a pod", func() {
			Eventually(func() ([]string, error) {
				return kubectl.GetOutput("logs", "windows-webserver")
			}, "30s").Should(Equal([]string{"Listening", "at", "http://*:80/"}))
		})

		expose := kubectl.StartKubectlCommand("expose", "pod", "windows-webserver", "--type", "NodePort")
		Eventually(expose, "30s").Should(gexec.Exit(0))

		By("should be able to reach it via NodePort", func() {
			hostIP := kubectl.GetOutputBytes("get", "pod", "-l", "app=windows-webserver",
				"-o", "jsonpath='{.items[0].status.hostIP}'")
			nodePort := kubectl.GetOutputBytes("get", "service", "windows-webserver",
				"-o", "jsonpath='{.spec.ports[0].nodePort}'")
			url := fmt.Sprintf("http://%s:%s", hostIP, nodePort)

			Eventually(curl(url), "30s").Should(ConsistOf("Windows", "Container", "Web", "Server"))
		})

		By("should be able to reach it via Cluster IP", func() {
			clusterIP := kubectl.GetOutputBytes("get", "service", "windows-webserver",
				"-o", "jsonpath='{.spec.clusterIP}'")
			url := fmt.Sprintf("http://%s", clusterIP)

			Eventually(curl(url), "30s").Should(ConsistOf("Windows", "Container", "Web", "Server"))
		})
	})
})

func curl(url string) func() ([]string, error) {
	name := fmt.Sprintf("curl-%d", rand.Int())
	job := fmt.Sprintf("job-name=%s", name)
	Eventually(
		kubectl.StartKubectlCommand("run", name, "--image=tutum/curl", "--restart=OnFailure",
			"--", "curl", "-s", url),
	).Should(gexec.Exit(0))

	Eventually(func() ([]string, error) {
		return kubectl.GetOutput("get", "pod", "-l", job, "-o", "jsonpath='{.items[0].status.phase}'")
	}, "30s").Should(ConsistOf("Succeeded"))

	return func() ([]string, error) {
		return kubectl.GetOutput("logs", "-l", job)
	}
}

func setupNS() {
	if !hasWindowsWorkers {
		Skip("skipping Windows tests since no Windows nodes were detected")
	}
	kubectl = test_helpers.NewKubectlRunner()
	kubectl.Setup()
}

func teardown() {
	kubectl.Teardown()
}
