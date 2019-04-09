package test_helpers

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/onsi/gomega/gexec"

	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type KubectlRunner struct {
	configPath       string
	namespace        string
	TimeoutInSeconds float64
}

func NewKubectlRunner() *KubectlRunner {
	return &KubectlRunner{
		namespace:        "test-" + GenerateRandomUUID(),
		TimeoutInSeconds: 60,
	}
}

func PathFromRoot(relativePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "..", "..", relativePath)
}

func (kubectl KubectlRunner) Setup() {
	kubectl.StartKubectlCommand("create", "namespace", kubectl.Namespace()).Wait(kubectl.TimeoutInSeconds)
}

func (kubectl KubectlRunner) Teardown() {
	kubectl.StartKubectlCommand("delete", "namespace", kubectl.Namespace()).Wait(kubectl.TimeoutInSeconds * 2)
}

func (kubectl KubectlRunner) Namespace() string {
	return kubectl.namespace
}

func (kubectl KubectlRunner) StartKubectlCommand(args ...string) *gexec.Session {
	return kubectl.StartKubectlCommandInNamespace(kubectl.namespace, args...)
}

func (kubectl KubectlRunner) StartKubectlCommandOrError(args ...string) (*gexec.Session, error) {
	return kubectl.StartKubectlCommandInNamespaceOrError(kubectl.namespace, args...)
}

func (kubectl KubectlRunner) RunKubectlCommandWithTimeout(args ...string) {
	Eventually(kubectl.StartKubectlCommandInNamespace(kubectl.namespace, args...), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
}

func (kubectl KubectlRunner) StartKubectlCommandInNamespace(namespace string, args ...string) *gexec.Session {
	argsWithNamespace := append([]string{"--namespace", namespace}, args...)
	if kubectl.configPath != "" {
		argsWithNamespace = append([]string{"--kubeconfig", kubectl.configPath}, argsWithNamespace...)
	}
	command := exec.Command("kubectl", argsWithNamespace...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session
}

func (kubectl KubectlRunner) StartKubectlCommandInNamespaceOrError(namespace string, args ...string) (*gexec.Session, error) {
	newArgs := append([]string{"--kubeconfig", kubectl.configPath, "--namespace", namespace}, args...)
	command := exec.Command("kubectl", newArgs...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	return session, err
}

func (kubectl KubectlRunner) StartKubectlCommandInNamespaceSilent(namespace string, args ...string) *gexec.Session {
	newArgs := append([]string{"--kubeconfig", kubectl.configPath, "--namespace", namespace}, args...)
	command := exec.Command("kubectl", newArgs...)
	session, err := gexec.Start(command, nil, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session
}

func GenerateRandomUUID() string {
	randomUUID := uuid.NewV4()
	return randomUUID.String()
}

func (kubectl *KubectlRunner) GetOutput(kubectlArgs ...string) ([]string, error) {
	output, err := kubectl.GetOutputBytesOrError(kubectlArgs...)
	return strings.Fields(string(output)), err
}

func (kubectl *KubectlRunner) GetOutputInNamespace(namespace string, kubectlArgs ...string) []string {
	output := kubectl.GetOutputBytesInNamespace(namespace, kubectlArgs...)
	return strings.Fields(string(output))
}

func (kubectl *KubectlRunner) GetOutputBytes(kubectlArgs ...string) []byte {
	session := kubectl.StartKubectlCommand(kubectlArgs...)
	Eventually(session, kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
	output := session.Out.Contents()
	return bytes.Trim(output, `"`)
}

func (kubectl *KubectlRunner) GetOutputBytesOrError(kubectlArgs ...string) ([]byte, error) {
	session := kubectl.StartKubectlCommand(kubectlArgs...)
	Eventually(session, kubectl.TimeoutInSeconds).Should(gexec.Exit())
	if session.ExitCode() != 0 {
		return []byte{}, fmt.Errorf("kubectl command exitted with non zero exit code: %d", session.ExitCode())
	}
	output := session.Out.Contents()
	return bytes.Trim(output, `"`), nil
}

func (kubectl *KubectlRunner) GetOutputBytesInNamespace(namespace string, kubectlArgs ...string) []byte {
	var session *gexec.Session
	Eventually(func() int {
		session = kubectl.StartKubectlCommandInNamespace(namespace, kubectlArgs...)
		Eventually(session, kubectl.TimeoutInSeconds).Should(gexec.Exit())

		return session.ExitCode()
	}, kubectl.TimeoutInSeconds, "30s").Should(Equal(0))
	output := session.Out.Contents()
	return bytes.Trim(output, `"`)
}

func (kubectl *KubectlRunner) GetNodePort(service string) (string, error) {
	return kubectl.GetNodePortInNamespace(service, kubectl.Namespace())
}

func (kubectl *KubectlRunner) GetNodePortInNamespace(service string, namespace string) (string, error) {
	output := kubectl.GetOutputInNamespace(namespace, "describe", service)

	for i := 0; i < len(output); i++ {
		if output[i] == "NodePort:" {
			nodePort := output[i+2]
			return nodePort[:strings.Index(nodePort, "/")], nil
		}
	}

	return "", errors.New("No nodePort found!")
}

func (kubectl *KubectlRunner) getWorkerIP() string {
	output, err := kubectl.GetOutput("get", "nodes", "-o", "jsonpath={.items[*].status.addresses[?(@.type==\"InternalIP\")].address}")
	Expect(err).NotTo(HaveOccurred())
	return output[0]
}

func (kubectl *KubectlRunner) GetAppAddress(service string) string {
	return kubectl.GetAppAddressInNamespace(service, kubectl.Namespace())
}

func (kubectl *KubectlRunner) GetAppAddressInNamespace(service string, namespace string) string {
	workerIP := kubectl.getWorkerIP()
	nodePort, err := kubectl.GetNodePortInNamespace(service, namespace)
	Expect(err).ToNot(HaveOccurred())

	return fmt.Sprintf("%s:%s", workerIP, nodePort)
}

func (kubectl *KubectlRunner) GetPodStatus(namespace string, podName string) string {
	return kubectl.getPodStatus(namespace, podName)
}

func (kubectl *KubectlRunner) GetResourceNameBySelector(namespace, resource, selector string) string {
	return kubectl.GetOutputInNamespace(namespace, "get", resource, "-l", selector, "-o", "jsonpath={.items[0].metadata.name}")[0]
}

func (kubectl *KubectlRunner) GetPodStatusBySelector(namespace string, selector string) string {
	return kubectl.getPodStatus(namespace, "-l", selector)
}

func (kubectl *KubectlRunner) getPodStatus(namespace string, selector ...string) string {
	args := []string{"describe", "pod"}
	args = append(args, selector...)
	output := kubectl.RunKubectlCommandWithRetry(namespace, kubectl.TimeoutInSeconds*2, args...)

	re := regexp.MustCompile(`Status:\s+(\w+)`)
	matches := re.FindStringSubmatch(output)
	podStatus := matches[1]
	return podStatus
}

// RunKubectlCommandWithRetry will run kubectl command with retries, until the timeout reaches
// the command will retry every 10s
// Expect the command output to be not empty
// Expect the command to exit 0
// return the command output
func (kubectl *KubectlRunner) RunKubectlCommandWithRetry(namespace string, timeout float64, args ...string) string {
	var session *gexec.Session

	Eventually(func() string {
		session = kubectl.StartKubectlCommandInNamespace(namespace, args...)
		Eventually(session, "10s").Should(gexec.Exit(0))

		return string(session.Out.Contents())
	}, timeout).ShouldNot(BeEmpty())

	return string(session.Out.Contents())
}

func (kubectl *KubectlRunner) GetLBAddress(service, iaas string) string {
	var jsonPathForLoadBalancer string

	if iaas == "gcp" || iaas == "gce" || iaas == "azure" { // TODO: remove GCP once testconfig is gone
		jsonPathForLoadBalancer = "jsonpath={.status.loadBalancer.ingress[0].ip}"
	} else if iaas == "aws" {
		jsonPathForLoadBalancer = "jsonpath={.status.loadBalancer.ingress[0].hostname}"
	}

	output, err := kubectl.GetOutput("get", "service", service, "-o", jsonPathForLoadBalancer)

	if err != nil {
		fmt.Fprintf(GinkgoWriter, "error when connecting to Kubernetes: %s", err.Error())
		return ""
	}

	if len(output) == 0 {
		fmt.Fprintf(GinkgoWriter, "loadbalancer still pending creation\n")
		return ""
	}

	fmt.Fprintf(GinkgoWriter, "Output %#v\n", output)
	if len(output) != 0 {
		return output[0]
	}
	return ""
}
