package test_helpers

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/onsi/gomega/gexec"

	testconfig "tests/config"

	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type KubectlRunner struct {
	configPath string
	namespace  string
	Timeout    string
}

func NewKubectlRunner(pathToKubeConfig string) *KubectlRunner {

	runner := &KubectlRunner{}

	runner.configPath = pathToKubeConfig
	if runner.configPath == "" {
		Fail("path to kubeconfig must be specified")
	}

	runner.namespace = "test-" + GenerateRandomUUID()
	runner.Timeout = "60s"

	return runner
}

func NewKubectlRunnerWithDefaultConfig() *KubectlRunner {
	return &KubectlRunner{
		namespace: "test-" + GenerateRandomUUID(),
		Timeout:   "60s",
	}
}

func PathFromRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "..", "..", relativePath)
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
		argsWithNamespace = append(argsWithNamespace, []string{"--kubeconfig", runner.configPath}...)
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

func (runner KubectlRunner) CreateNamespace() {
	Eventually(runner.RunKubectlCommand("create", "namespace", runner.namespace), "60s").Should(gexec.Exit(0))
}

func (runner *KubectlRunner) GetOutput(kubectlArgs ...string) []string {
	output := runner.GetOutputBytes(kubectlArgs...)
	return strings.Fields(string(output))
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
	session := runner.RunKubectlCommandInNamespace(namespace, kubectlArgs...)
	Eventually(session, "60s").Should(gexec.Exit(0))
	output := session.Out.Contents()
	return bytes.Trim(output, `"`)
}

func (runner *KubectlRunner) GetNodePort(service string) (string, error) {
	output := runner.GetOutput("describe", service)

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
	output := runner.GetOutput("get", "nodes", "-o", "jsonpath=\"{.items[0].metadata.labels['spec\\.ip']}\"")
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

func (runner *KubectlRunner) GetServiceAccount(deployment, namespace string) string {
	s := runner.RunKubectlCommandInNamespace(namespace, "get", "deployment/"+deployment,
		"-o", "jsonpath='{.spec.template.spec.serviceAccountName}'")
	Eventually(s, "15s").Should(gexec.Exit(0))
	return string(s.Out.Contents())
}

func (runner *KubectlRunner) GetPodStatus(namespace string, podName string) string {
	session := runner.RunKubectlCommandInNamespace(namespace, "describe", "pod", podName)
	Eventually(session, "120s").Should(gexec.Exit(0))
	re := regexp.MustCompile(`Status:\s+(\w+)`)
	matches := re.FindStringSubmatch(string(session.Out.Contents()))
	podStatus := matches[1]
	return podStatus
}

func (runner *KubectlRunner) GetLBAddress(service, iaas string) string {
	output := []string{}
	loadBalancerAddress := ""
	if iaas == "gcp" {
		output = runner.GetOutput("get", "service", service, "-o", "jsonpath={.status.loadBalancer.ingress[0].ip}")
	} else if iaas == "aws" {
		output = runner.GetOutput("get", "service", service, "-o", "jsonpath={.status.loadBalancer.ingress[0].hostname}")
	}

	if len(output) == 0 {
		fmt.Printf("loadbalancer still pending creation\n")
		return ""
	}

	fmt.Printf("Output %#v", output)
	if len(output) != 0 {
		loadBalancerAddress = output[0]
	}
	return loadBalancerAddress
}

func (runner *KubectlRunner) CleanupServiceWithLB(loadBalancerAddress, pathToSpec, iaas string, aws testconfig.AWS) {
	lbSecurityGroup := ""

	if iaas == "aws" {
		// Get the LB
		if loadBalancerAddress != "" {
			// Get the security group
			cmd := exec.Command("aws", "elb", "describe-load-balancers", "--region", aws.Region, "--query",
				fmt.Sprintf("LoadBalancerDescriptions[?DNSName==`%s`].[SecurityGroups]", loadBalancerAddress),
				"--output", "text")
			cmd.Env = append(os.Environ(),
				fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", aws.AccessKeyID),
				fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", aws.SecretAccessKey),
			)
			fmt.Fprintf(GinkgoWriter, "Get LoadBalancer security group - %s\n", cmd.Args)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Eventually(session, "10s").Should(gexec.Exit(0))
			Expect(err).NotTo(HaveOccurred())
			output := strings.Fields(string(session.Out.Contents()))
			if len(output) != 0 {
				lbSecurityGroup = output[0]
				fmt.Printf("Found LB security group [%s]", lbSecurityGroup)
			}

		}
	}

	session := runner.RunKubectlCommand("delete", "-f", pathToSpec)
	session.Wait("60s")

	// Teardown the security group
	if lbSecurityGroup != "" {
		cmd := exec.Command("aws", "ec2", "revoke-security-group-ingress", "--region", aws.Region, "--group-id",
			aws.IngressGroupID, "--source-group", lbSecurityGroup, "--protocol", "all")
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", aws.AccessKeyID),
			fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", aws.SecretAccessKey),
		)

		fmt.Fprintf(GinkgoWriter, "Teardown security groups - %s\n", cmd.Args)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session, "10s").Should(gexec.Exit())
		fmt.Printf("Tearing down security group exited with code '%d'\n", session.ExitCode())
	}
}
