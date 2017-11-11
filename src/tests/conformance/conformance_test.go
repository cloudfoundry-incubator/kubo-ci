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
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Conformance Tests", func() {
	var conformanceSpec string
	var kubectl *KubectlRunner

	BeforeEach(func() {
		conformanceSpec = GetLatestConformanceSpec()
		kubectl = NewKubectlRunner()
	})

	AfterEach(func() {
		kubectl.RunKubectlCommandInNamespace("sonobuoy", "delete", "-f", conformanceSpec)
		os.Remove(conformanceSpec)
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

		By("Extracting the conformance test results")
		session = kubectl.RunKubectlCommandInNamespace("sonobuoy", "log", "sonobuoy")
		re := regexp.MustCompile(`/tmp/sonobuoy/.*\.tar.gz`)
		matches := re.FindStringSubmatch(string(session.Out.Contents()))
		Expect(len(matches)).To(Equal(1))
		logPath := matches[0]
		containerAddressedLogPath := fmt.Sprintf("sonobuoy:%s", logPath)
		conformanceResultsPath := os.Getenv("CONFORMANCE_RESULTS_PATH")
		session = kubectl.RunKubectlCommandInNamespace("sonobuoy", "cp", containerAddressedLogPath, conformanceResultsPath)
		Eventually(session, "60s").Should(gexec.Exit(0))
		dir, err := ioutil.TempDir("", "results")
		Expect(err).NotTo(HaveOccurred())
		command := exec.Command("tar", "xvf", conformanceResultsPath, "-C", dir)
		err = command.Start()
		Expect(err).NotTo(HaveOccurred())

		By("Reading the test results")
		e2eLogPath := filepath.Join(dir, "plugins/e2e/results/e2e.log")
		re = regexp.MustCompile(`(\d+) Failed .* TestE2E`)
		e2eLogContents, err := ioutil.ReadFile(e2eLogPath)
		Expect(err).NotTo(HaveOccurred())

		matches = re.FindStringSubmatch(string(e2eLogContents))
		Expect(len(matches)).To(Equal(2))

		numFailures := matches[1]
		Expect(numFailures).To(Equal("0"))
	})
})

func GetLatestConformanceSpec() string {
	resp, err := http.Get("https://raw.githubusercontent.com/cncf/k8s-conformance/master/sonobuoy-conformance.yaml")
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
