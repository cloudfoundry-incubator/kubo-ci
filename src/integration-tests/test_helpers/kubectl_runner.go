package test_helpers

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type KubectlRunner struct {
	configPath string
	namespace string
}

func NewKubectlRunner() *KubectlRunner {

	runner := &KubectlRunner{}

	runner.configPath = os.Getenv("PATH_TO_KUBECONFIG")
	if runner.configPath == "" {
		Fail("PATH_TO_KUBECONFIG is not set")
	}

	runner.namespace = "test-" + GenerateRandomName()

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
	newArgs := append([]string{"--kubeconfig", runner.configPath, "--namespace", runner.namespace}, args...)
	command := exec.Command("kubectl", newArgs...)

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session

}

func GenerateRandomName() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 20)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}