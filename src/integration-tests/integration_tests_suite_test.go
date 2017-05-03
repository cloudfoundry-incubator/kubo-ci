package integration_tests_test

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

var (
	pathToKubeConfig string
	kubeNamespace    string
	appsDomain       string
	tcpRouterDNSName string
	tcpPort          int64
)

var _ = BeforeSuite(func() {
	var err error

	pathToKubeConfig = os.Getenv("PATH_TO_KUBECONFIG")
	if pathToKubeConfig == "" {
		Fail("PATH_TO_KUBECONFIG is not set")
	}

	tcpPort, err = strconv.ParseInt(os.Getenv("WORKLOAD_TCP_PORT"), 10, 64)
	if err != nil || tcpPort <= 0 {
		Fail("Correct WORKLOAD_TCP_PORT has to be set")
	}

	appsDomain = os.Getenv("CF_APPS_DOMAIN")
	if appsDomain == "" {
		Fail("CF_APPS_DOMAIN is not set")
	}

	tcpRouterDNSName = os.Getenv("TCP_ROUTER_DNS_NAME")
	if tcpRouterDNSName == "" {
		Fail("TCP_ROUTER_DNS_NAME is not set")
	}

	kubeNamespace = "test-" + generateRandomName()
	runKubectlCommand("create", "namespace", kubeNamespace).Wait("60s")
})

var _ = AfterSuite(func() {
	if kubeNamespace != "" {
		runKubectlCommand("delete", "namespace", kubeNamespace).Wait("60s")
	}
})

func runKubectlCommand(args ...string) *gexec.Session {
	newArgs := append([]string{"--kubeconfig", pathToKubeConfig, "--namespace", kubeNamespace}, args...)
	command := exec.Command("kubectl", newArgs...)

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session

}

func pathFromRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "..", relativePath)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandomName() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 20)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
