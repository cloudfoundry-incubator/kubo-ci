package conformance_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"tests/config"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type Manifest struct {
	Version string `yaml:"version"`
}

var _ = Describe("Conformance Tests", func() {
	var (
		conformanceSpec string
		kubectl         *KubectlRunner
		testconfig      *config.Config
	)

	BeforeSuite(func() {
		var err error
		testconfig, err = config.InitConfig()
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		conformanceSpec = GetLatestConformanceSpec()
		kubectl = NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
	})

	AfterEach(func() {
		if !CurrentGinkgoTestDescription().Failed {
			session := kubectl.RunKubectlCommandInNamespace("sonobuoy", "delete", "-f", conformanceSpec)
			Eventually(session, "30s").Should(gexec.Exit(0))
			os.Remove(conformanceSpec)
		}
	})

	Specify("Conformance tests succeeds", func() {

		By("Applying the conformance spec")
		session := kubectl.RunKubectlCommandInNamespace("sonobuoy", "apply", "-f", conformanceSpec)
		Eventually(session, "30s").Should(gexec.Exit(0))

		By("Waiting for sonobuoy pod to be running")
		Eventually(func() string {
			outputs := kubectl.GetOutput("get", "pod/sonobuoy", "-n", "sonobuoy", "-o", "jsonpath={.status.phase}")
			return string(outputs[0])
		}, "60s", "2s").Should(Equal("Running"))

		By("Waiting for conformance tests to complete")
		Eventually(func() string {
			outputs := kubectl.GetOutput("log", "sonobuoy", "-n", "sonobuoy")
			return strings.Join(outputs, " ")
		}, "60m", "1m").Should(ContainSubstring("no-exit was specified, sonobuoy is now blocking"))

		By("Locate test results")
		session = kubectl.RunKubectlCommandInNamespace("sonobuoy", "log", "sonobuoy")
		Eventually(session, "20s").Should(gexec.Exit(0))
		re := regexp.MustCompile(`/tmp/sonobuoy/.*\.tar.gz`)

		conformanceLogs := string(session.Out.Contents())
		fmt.Println("Grabbing logs tarball...")
		fmt.Println(conformanceLogs)
		matches := re.FindStringSubmatch(conformanceLogs)
		Expect(len(matches)).To(Equal(1))
		logPath := matches[0]

		By("Get the release version")
		releaseVersion := os.Getenv("CONFORMANCE_RELEASE_VERSION")
		fmt.Println(fmt.Sprintf("release version: %s", releaseVersion))

		By("Move results to output dir")
		conformanceResultsDir := os.Getenv("CONFORMANCE_RESULTS_DIR")
		fmt.Println(fmt.Sprintf("conformance results dir: %s", conformanceResultsDir))
		conformanceResultsPath := filepath.Join(conformanceResultsDir, fmt.Sprintf("conformance-results-%s.tar.gz", releaseVersion))
		containerAddressedLogPath := fmt.Sprintf("sonobuoy:%s", logPath)
		session = kubectl.RunKubectlCommandInNamespace("sonobuoy", "cp", containerAddressedLogPath, conformanceResultsPath)
		Eventually(session, "60s").Should(gexec.Exit(0))
		dir, err := ioutil.TempDir("", "results")
		Expect(err).NotTo(HaveOccurred())

		By("Extract results")
		command := exec.Command("tar", "xvf", conformanceResultsPath, "-C", dir)
		err = command.Run()
		Expect(err).NotTo(HaveOccurred())

		By("Reading the test results")
		e2eLogPath := filepath.Join(dir, "plugins/e2e/results/e2e.log")
		re = regexp.MustCompile(`(FAIL|SUCCESS)! -- (\d+) Passed \| (\d+) Failed \| (\d+) Pending \| (\d+) Skipped`)
		e2eLogContents, err := ioutil.ReadFile(e2eLogPath)
		Expect(err).NotTo(HaveOccurred())
		fmt.Println("E2E Test log:")
		fmt.Println(string(e2eLogContents))

		matches = re.FindStringSubmatch(string(e2eLogContents))
		Expect(len(matches)).To(Equal(6))

		testSuiteResult := matches[1]
		Expect(testSuiteResult).To(Equal("SUCCESS"))

		numFailures := matches[3]
		Expect(numFailures).To(Equal("0"))
	})
})

func GetLatestConformanceSpec() string {
	resp, err := http.Get("https://raw.githubusercontent.com/cloudfoundry-incubator/k8s-conformance/e2e-fix/sonobuoy-conformance.yaml")
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	conformanceYaml, err := ioutil.TempFile("", "conformance")
	Expect(err).NotTo(HaveOccurred())
	contents, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	_, err = conformanceYaml.Write(contents)
	Expect(err).NotTo(HaveOccurred())
	return conformanceYaml.Name()
}
