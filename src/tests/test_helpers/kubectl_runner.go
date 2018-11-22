package test_helpers

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
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
	kubectl.RunKubectlCommand("create", "namespace", kubectl.Namespace()).Wait(kubectl.TimeoutInSeconds)
}

func (kubectl KubectlRunner) Teardown() {
	kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace()).Wait(kubectl.TimeoutInSeconds * 2)
}

func (kubectl KubectlRunner) Namespace() string {
	return kubectl.namespace
}

func (kubectl KubectlRunner) RunKubectlCommand(args ...string) *gexec.Session {
	return kubectl.RunKubectlCommandInNamespace(kubectl.namespace, args...)
}

func (kubectl KubectlRunner) RunKubectlCommandOrError(args ...string) (*gexec.Session, error) {
	return kubectl.RunKubectlCommandInNamespaceOrError(kubectl.namespace, args...)
}

func (kubectl KubectlRunner) RunKubectlCommandWithTimeout(args ...string) {
	Eventually(kubectl.RunKubectlCommandInNamespace(kubectl.namespace, args...), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
}

func (kubectl KubectlRunner) RunKubectlCommandInNamespace(namespace string, args ...string) *gexec.Session {
	argsWithNamespace := append([]string{"--namespace", namespace}, args...)
	if kubectl.configPath != "" {
		argsWithNamespace = append([]string{"--kubeconfig", kubectl.configPath}, argsWithNamespace...)
	}
	command := exec.Command("kubectl", argsWithNamespace...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session
}

func (kubectl KubectlRunner) RunKubectlCommandInNamespaceOrError(namespace string, args ...string) (*gexec.Session, error) {
	newArgs := append([]string{"--kubeconfig", kubectl.configPath, "--namespace", namespace}, args...)
	command := exec.Command("kubectl", newArgs...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	return session, err
}

func (kubectl KubectlRunner) RunKubectlCommandInNamespaceSilent(namespace string, args ...string) *gexec.Session {
	newArgs := append([]string{"--kubeconfig", kubectl.configPath, "--namespace", namespace}, args...)
	command := exec.Command("kubectl", newArgs...)
	session, err := gexec.Start(command, nil, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session
}

func (kubectl KubectlRunner) ExpectEventualSuccess(args ...string) {
	Eventually(kubectl.RunKubectlCommand(args...), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
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
	session := kubectl.RunKubectlCommand(kubectlArgs...)
	Eventually(session, kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
	output := session.Out.Contents()
	return bytes.Trim(output, `"`)
}

func (kubectl *KubectlRunner) GetOutputBytesOrError(kubectlArgs ...string) ([]byte, error) {
	session := kubectl.RunKubectlCommand(kubectlArgs...)
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
		session = kubectl.RunKubectlCommandInNamespace(namespace, kubectlArgs...)
		Eventually(session, kubectl.TimeoutInSeconds).Should(gexec.Exit())

		return session.ExitCode()
	}, kubectl.TimeoutInSeconds, "30s").Should(Equal(0))
	output := session.Out.Contents()
	return bytes.Trim(output, `"`)
}

func (kubectl *KubectlRunner) GetNodePort(service string) (string, error) {
	output, err := kubectl.GetOutput("describe", service)
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

func (kubectl *KubectlRunner) GetWorkerIP() string {
	output, err := kubectl.GetOutput("get", "nodes", "-o", "jsonpath={.items[*].status.addresses[?(@.type==\"InternalIP\")].address}")
	Expect(err).NotTo(HaveOccurred())
	return output[0]
}

func (kubectl *KubectlRunner) GetAppAddress(service string) string {
	workerIP := kubectl.GetWorkerIP()
	nodePort, err := kubectl.GetNodePort(service)
	Expect(err).ToNot(HaveOccurred())

	return fmt.Sprintf("%s:%s", workerIP, nodePort)
}

func (kubectl *KubectlRunner) GetAppAddressInNamespace(service string, namespace string) string {
	workerIP := kubectl.GetWorkerIP()
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
	var session *gexec.Session
	args := []string{"describe", "pod"}
	args = append(args, selector...)
	Eventually(func() string {
		session = kubectl.RunKubectlCommandInNamespace(namespace, args...)
		Eventually(session, "10s").Should(gexec.Exit(0))

		return string(session.Out.Contents())
	}, kubectl.TimeoutInSeconds*2).ShouldNot(BeEmpty())

	re := regexp.MustCompile(`Status:\s+(\w+)`)
	matches := re.FindStringSubmatch(string(session.Out.Contents()))
	podStatus := matches[1]
	return podStatus
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

func (kubectl *KubectlRunner) templatePSPWithNamespace(namespace string) string {
	_, filename, _, _ := runtime.Caller(0)
	srcDir, err := filepath.Abs(filepath.Dir(filename))
	Expect(err).NotTo(HaveOccurred())

	file := filepath.Join(srcDir, "fixtures", "smoke-test-psp.yml")

	t, err := template.ParseFiles(file)
	Expect(err).NotTo(HaveOccurred())

	f, err := ioutil.TempFile("", filepath.Base(file))
	Expect(err).NotTo(HaveOccurred())
	defer f.Close()

	type templateInfo struct{ PSPName, Namespace string }
	Expect(t.Execute(f, templateInfo{PSPName: namespace, Namespace: namespace})).To(Succeed())

	return f.Name()
}
