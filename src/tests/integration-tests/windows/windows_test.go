package windows_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	kubectl         *test_helpers.KubectlRunner
	webServerSpec   = test_helpers.PathFromRoot("specs/windows/webserver.yml")
	curlWindowsSpec = test_helpers.PathFromRoot("specs/windows/curl.yml")
	curlPod         = v1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		Spec: v1.PodSpec{
			NodeSelector: map[string]string{"beta.kubernetes.io/os": "windows"},
			Tolerations: []v1.Toleration{
				{Key: "windows", Operator: "Equal", Effect: "NoSchedule", Value: "2019"},
			},
			RestartPolicy: v1.RestartPolicyNever,
			Containers: []v1.Container{
				{Name: "curl", Image: "mcr.microsoft.com/windows/nanoserver:1809", Command: []string{"curl.exe"}},
			},
		},
	}
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
			}, "30s").Should(Equal([]string{"Proudly", "serving", "content", "on", "port", "80"}))
		})

		expose := kubectl.StartKubectlCommand("expose", "pod", "windows-webserver", "--type", "NodePort")
		Eventually(expose, "30s").Should(gexec.Exit(0))

		By("should be able to reach it via NodePort", func() {
			hostIP := kubectl.GetOutputBytes("get", "pod", "-l", "app=windows-webserver",
				"-o", "jsonpath='{.items[0].status.hostIP}'")
			nodePort := kubectl.GetOutputBytes("get", "service", "windows-webserver",
				"-o", "jsonpath='{.spec.ports[0].nodePort}'")
			url := fmt.Sprintf("http://%s:%s", hostIP, nodePort)

			Eventually(curlLinux(url), "30s").Should(ContainElement(ContainSubstring("webserver.exe")))
			Eventually(curlWindows(url), "30s").Should(ContainElement(ContainSubstring("webserver.exe")))
		})

		By("should be able to reach it via Cluster IP", func() {
			clusterIP := kubectl.GetOutputBytes("get", "service", "windows-webserver",
				"-o", "jsonpath='{.spec.clusterIP}'")
			url := fmt.Sprintf("http://%s", clusterIP)

			Eventually(curlLinux(url), "100s").Should(ContainElement(ContainSubstring("webserver.exe")))
			Eventually(curlWindows(url), "100s").Should(ContainElement(ContainSubstring("webserver.exe")))
		})
	})
})

func curlLinux(url string) func() ([]string, error) {
	name := fmt.Sprintf("curl-%d", rand.Int())
	Eventually(
		kubectl.StartKubectlCommand("run", name, "--image=tutum/curl", "--restart=Never",
			"--", "curl", "-s", url),
	).Should(gexec.Exit(0))

	Eventually(func() ([]string, error) {
		return kubectl.GetOutput("get", "pod", name, "-o", "jsonpath='{.status.phase}'")
	}, "30s").Should(ConsistOf("Succeeded"))

	return func() ([]string, error) {
		return kubectl.GetOutput("logs", name)
	}
}

func curlWindows(url string) func() ([]string, error) {
	name := fmt.Sprintf("curl-windows-%d", rand.Int())

	curlPod.Spec.Containers[0].Args = []string{"-s", url}
	curlPod.Name = name

	outSpec, err := ioutil.TempFile("", "curl")
	Expect(err).To(BeNil())
	defer os.Remove(outSpec.Name())

	Expect(json.NewEncoder(outSpec).Encode(&curlPod)).To(Succeed())

	Eventually(
		kubectl.StartKubectlCommand("create", "-f", outSpec.Name()),
	).Should(gexec.Exit(0))

	Eventually(func() ([]string, error) {
		return kubectl.GetOutput("get", "pod", name, "-o", "jsonpath='{.status.phase}'")
	}, "60s").Should(ConsistOf("Succeeded"))

	return func() ([]string, error) {
		return kubectl.GetOutput("logs", name)
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
