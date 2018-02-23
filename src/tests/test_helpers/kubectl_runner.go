package test_helpers

import (
	"errors"
	"math/rand"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"strings"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/gomega/gexec"
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

	runner.namespace = "test-" + GenerateRandomName()
	runner.Timeout = "60s"

	return runner
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

func (runner KubectlRunner) RunKubectlCommandWithTimeout(args ...string) {
	Eventually(runner.RunKubectlCommandInNamespace(runner.namespace, args...), "60s").Should(gexec.Exit(0))
}

func (runner KubectlRunner) RunKubectlCommandInNamespace(namespace string, args ...string) *gexec.Session {
	newArgs := append([]string{"--kubeconfig", runner.configPath, "--namespace", namespace}, args...)
	command := exec.Command("kubectl", newArgs...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session
}

func (runner KubectlRunner) ExpectEventualSuccess(args ...string) {
	Eventually(runner.RunKubectlCommand(args...), runner.Timeout).Should(gexec.Exit(0))
}

func GenerateRandomName() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 20)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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
	Eventually(session, "20s").Should(gexec.Exit(0))
	output := session.Out.Contents()
	return output
}

func (runner *KubectlRunner) GetOutputBytesInNamespace(namespace string, kubectlArgs ...string) []byte {
	session := runner.RunKubectlCommandInNamespace(namespace, kubectlArgs...)
	Eventually(session, "20s").Should(gexec.Exit(0))
	output := session.Out.Contents()
	return output
}

func init() {
	rand.Seed(config.GinkgoConfig.RandomSeed)
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

func (runner *KubectlRunner) GetAppAddress(deployment director.Deployment, service string) string {
	workerIP := GetWorkerIP(deployment)
	nodePort, err := runner.GetNodePort(service)
	Expect(err).ToNot(HaveOccurred())

	return fmt.Sprintf("%s:%s", workerIP, nodePort)
}

func (runner *KubectlRunner) GetAppAddressInNamespace(deployment director.Deployment, service string, namespace string) string {
	workerIP := GetWorkerIP(deployment)
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
	fmt.Printf("Output [%s]", output)
	if len(output) != 0 {
		loadBalancerAddress = output[0]
	}
	return loadBalancerAddress
}

func (runner *KubectlRunner) CleanupServiceWithLB(loadBalancerAddress, pathToSpec, iaas string) {
	lbSecurityGroup := ""

	if iaas == "aws" {
		// Get the LB
		if loadBalancerAddress != "" {
			// Get the security group
			cmd := exec.Command("aws", "elb", "describe-load-balancers", "--query",
				fmt.Sprintf("LoadBalancerDescriptions[?DNSName==`%s`].[SecurityGroups]", loadBalancerAddress),
				"--output", "text")
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
		cmd := exec.Command("aws", "ec2", "revoke-security-group-ingress", "--group-id",
			os.Getenv("AWS_INGRESS_GROUP_ID"), "--source-group", lbSecurityGroup, "--protocol", "all")
		fmt.Fprintf(GinkgoWriter, "Teardown security groups - %s\n", cmd.Args)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session, "10s").Should(gexec.Exit(0))
	}
}
