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
	configPath string
	namespace  string
	Timeout    string
}

func NewKubectlRunner() *KubectlRunner {
	return &KubectlRunner{
		namespace: "test-" + GenerateRandomUUID(),
		Timeout:   "60s",
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

func (runner KubectlRunner) Setup() {
	runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
}

func (runner KubectlRunner) Teardown() {
	runner.RunKubectlCommand("delete", "namespace", runner.Namespace())
}

func (runner KubectlRunner) Namespace() string {
	return runner.namespace
}

func (runner KubectlRunner) RunKubectlCommand(args ...string) *gexec.Session {
	return runner.RunKubectlCommandInNamespace(runner.namespace, args...)
}

func (runner KubectlRunner) RunKubectlCommandOrError(args ...string) (*gexec.Session, error) {
	return runner.RunKubectlCommandInNamespaceOrError(runner.namespace, args...)
}

func (runner KubectlRunner) RunKubectlCommandWithTimeout(args ...string) {
	Eventually(runner.RunKubectlCommandInNamespace(runner.namespace, args...), "60s").Should(gexec.Exit(0))
}

func (runner KubectlRunner) RunKubectlCommandInNamespace(namespace string, args ...string) *gexec.Session {
	argsWithNamespace := append([]string{"--namespace", namespace}, args...)
	if runner.configPath != "" {
		argsWithNamespace = append([]string{"--kubeconfig", runner.configPath}, argsWithNamespace...)
	}
	command := exec.Command("kubectl", argsWithNamespace...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session
}

func (runner KubectlRunner) RunKubectlCommandInNamespaceOrError(namespace string, args ...string) (*gexec.Session, error) {
	newArgs := append([]string{"--kubeconfig", runner.configPath, "--namespace", namespace}, args...)
	command := exec.Command("kubectl", newArgs...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	return session, err
}

func (runner KubectlRunner) RunKubectlCommandInNamespaceSilent(namespace string, args ...string) *gexec.Session {
	newArgs := append([]string{"--kubeconfig", runner.configPath, "--namespace", namespace}, args...)
	command := exec.Command("kubectl", newArgs...)
	session, err := gexec.Start(command, nil, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session
}

func (runner KubectlRunner) ExpectEventualSuccess(args ...string) {
	Eventually(runner.RunKubectlCommand(args...), runner.Timeout).Should(gexec.Exit(0))
}

func GenerateRandomUUID() string {
	randomUUID := uuid.NewV4()
	return randomUUID.String()
}

func (runner *KubectlRunner) GetOutput(kubectlArgs ...string) ([]string, error) {
	output, err := runner.GetOutputBytesOrError(kubectlArgs...)
	return strings.Fields(string(output)), err
}

func (runner *KubectlRunner) GetOutputInNamespace(namespace string, kubectlArgs ...string) []string {
	output := runner.GetOutputBytesInNamespace(namespace, kubectlArgs...)
	return strings.Fields(string(output))
}

func (runner *KubectlRunner) GetOutputBytes(kubectlArgs ...string) []byte {
	session := runner.RunKubectlCommand(kubectlArgs...)
	Eventually(session, "60s").Should(gexec.Exit(0))
	output := session.Out.Contents()
	return bytes.Trim(output, `"`)
}

func (runner *KubectlRunner) GetOutputBytesOrError(kubectlArgs ...string) ([]byte, error) {
	session := runner.RunKubectlCommand(kubectlArgs...)
	Eventually(session, "60s").Should(gexec.Exit())
	if session.ExitCode() != 0 {
		return []byte{}, fmt.Errorf("kubectl command exitted with non zero exit code: %d", session.ExitCode())
	}
	output := session.Out.Contents()
	return bytes.Trim(output, `"`), nil
}

func (runner *KubectlRunner) GetOutputBytesInNamespace(namespace string, kubectlArgs ...string) []byte {
	var session *gexec.Session
	Eventually(func() int {
		session = runner.RunKubectlCommandInNamespace(namespace, kubectlArgs...)
		Eventually(session, "60s").Should(gexec.Exit())

		return session.ExitCode()
	}, "60s", "30s").Should(Equal(0))
	output := session.Out.Contents()
	return bytes.Trim(output, `"`)
}

func (runner *KubectlRunner) GetNodePort(service string) (string, error) {
	output, err := runner.GetOutput("describe", service)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(output); i++ {
		if output[i] == "NodePort:" {
			nodePort := output[i+2]
			return nodePort[:strings.Index(nodePort, "/")], nil
		}
	}

	return "", errors.New("No nodePort found!")
}

func (runner *KubectlRunner) GetNodePortInNamespace(service string, namespace string) (string, error) {
	output := runner.GetOutputInNamespace(namespace, "describe", service)

	for i := 0; i < len(output); i++ {
		if output[i] == "NodePort:" {
			nodePort := output[i+2]
			return nodePort[:strings.Index(nodePort, "/")], nil
		}
	}

	return "", errors.New("No nodePort found!")
}

func (runner *KubectlRunner) GetWorkerIP() string {
	output, err := runner.GetOutput("get", "nodes", "-o", "jsonpath={.items[*].status.addresses[?(@.type==\"InternalIP\")].address}")
	Expect(err).NotTo(HaveOccurred())
	return output[0]
}

func (runner *KubectlRunner) GetAppAddress(service string) string {
	workerIP := runner.GetWorkerIP()
	nodePort, err := runner.GetNodePort(service)
	Expect(err).ToNot(HaveOccurred())

	return fmt.Sprintf("%s:%s", workerIP, nodePort)
}

func (runner *KubectlRunner) GetAppAddressInNamespace(service string, namespace string) string {
	workerIP := runner.GetWorkerIP()
	nodePort, err := runner.GetNodePortInNamespace(service, namespace)
	Expect(err).ToNot(HaveOccurred())

	return fmt.Sprintf("%s:%s", workerIP, nodePort)
}

func (runner *KubectlRunner) GetPodStatus(namespace string, podName string) string {
	session := runner.RunKubectlCommandInNamespace(namespace, "describe", "pod", podName)
	Eventually(session, "120s").Should(gexec.Exit(0))
	re := regexp.MustCompile(`Status:\s+(\w+)`)
	matches := re.FindStringSubmatch(string(session.Out.Contents()))
	podStatus := matches[1]
	return podStatus
}

func (runner *KubectlRunner) GetResourceNameBySelector(namespace, resource, selector string) string {
	return runner.GetOutputInNamespace(namespace, "get", resource, "-l", selector, "-o", "jsonpath={.items[0].metadata.name}")[0]
}

func (runner *KubectlRunner) GetPodStatusBySelector(namespace string, selector string) string {
	var session *gexec.Session
	Eventually(func() string {
		session = runner.RunKubectlCommandInNamespace(namespace, "describe", "pod", "-l", selector)
		Eventually(session, "10s").Should(gexec.Exit(0))

		return string(session.Out.Contents())
	}, "120s").ShouldNot(BeEmpty())

	re := regexp.MustCompile(`Status:\s+(\w+)`)
	matches := re.FindStringSubmatch(string(session.Out.Contents()))
	podStatus := matches[1]
	return podStatus
}

func (runner *KubectlRunner) GetLBAddress(service, iaas string) string {
	var jsonPathForLoadBalancer string

	if iaas == "gcp" || iaas == "gce" || iaas == "azure" { // TODO: remove GCP once testconfig is gone
		jsonPathForLoadBalancer = "jsonpath={.status.loadBalancer.ingress[0].ip}"
	} else if iaas == "aws" {
		jsonPathForLoadBalancer = "jsonpath={.status.loadBalancer.ingress[0].hostname}"
	}

	output, err := runner.GetOutput("get", "service", service, "-o", jsonPathForLoadBalancer)

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
