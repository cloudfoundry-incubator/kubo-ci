package test_helpers

import (
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"strings"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/gomega/gexec"
)

type KubectlRunner struct {
	configPath string
	namespace  string
	Timeout    string
}

func NewKubectlRunner() *KubectlRunner {

	runner := &KubectlRunner{}

	runner.configPath = os.Getenv("PATH_TO_KUBECONFIG")
	if runner.configPath == "" {
		Fail("PATH_TO_KUBECONFIG is not set")
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
	Eventually(runner.RunKubectlCommand("create", "namespace", runner.namespace), "20s").Should(gexec.Exit(0))
}

func (runner *KubectlRunner) GetOutput(kubectlArgs ...string) []string {
	session := runner.RunKubectlCommand(kubectlArgs...)
	Eventually(session, "10s").Should(gexec.Exit(0))
	output := session.Out.Contents()

	return strings.Fields(string(output))
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

func (runner *KubectlRunner) GetAppAddress(deployment director.Deployment, service string) string {
	workerIP := GetWorkerIP(deployment)
	nodePort, err := runner.GetNodePort(service)
	Expect(err).ToNot(HaveOccurred())

	return fmt.Sprintf("%s:%s", workerIP, nodePort)
}
