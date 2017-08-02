package test_helpers

import (
	"os"
	"fmt"
	"path/filepath"
	"runtime"
	"bytes"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
)

type TurbulenceRunner struct {
	username    string
  password    string
	apiEndpoint string
}

func NewTurbulenceRunner() *TurbulenceRunner {

	runner := &TurbulenceRunner{}

	runner.apiEndpoint = os.Getenv("TURBULENCE_API_ENDPOINT")
	if runner.apiEndpoint == "" {
		Fail("TURBULENCE_API_ENDPOINT is not set")
	}

	runner.username = os.Getenv("TURBULENCE_USERNAME")
	if runner.username == "" {
		Fail("TURBULENCE_USERNAME is not set")
	}

	runner.password = os.Getenv("TURBULENCE_PASSWORD")
	if runner.password == "" {
		Fail("TURBULENCE_PASSWORD is not set")
	}

	return runner
}

func PathFromRoot(relativePath string) string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	return filepath.Join(currentDir, "..", "..", "..", relativePath)
}

func (runner TurbulenceRunner) ApplyIncident(incidentFile string) (string, error) {
	postUrl := fmt.Sprintf("https://%s:%s@%s/incidents", runner.username, runner.password, runner.apiEndpoint)
	fileBytes, fileReadErr := ioutil.ReadFile(incidentFile)
	if fileReadErr != nil {
		return "", fmt.Errorf("Can't read file [%s]; [%v]", incidentFile, fileReadErr)
	}

	response, postErr := http.Post(postUrl, "application/json", bytes.NewReader(fileBytes))
  if postErr != nil {
		return "", fmt.Errorf("Error submitting incident; [%v]", postErr)
	}

	bytes, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		return "", fmt.Errorf("Error parsing incident submit response; [%v]", bodyErr)
	}

	return string(bytes), nil
}


//
// func (runner TurbulenceRunner) RunKubectlCommandInNamespace(namespace string, args ...string) *gexec.Session {
// 	newArgs := append([]string{"--kubeconfig", runner.configPath, "--namespace", namespace}, args...)
// 	command := exec.Command("kubectl", newArgs...)
//
// 	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
//
// 	Expect(err).NotTo(HaveOccurred())
// 	return session
// }
//
// func (runner TurbulenceRunner) ExpectEventualSuccess(args ...string) {
// 	Eventually(runner.RunKubectlCommand(args...), runner.Timeout).Should(gexec.Exit(0))
// }
//
// func GenerateRandomName() string {
// 	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
// 	b := make([]rune, 20)
// 	for i := range b {
// 		b[i] = letterRunes[rand.Intn(len(letterRunes))]
// 	}
// 	return string(b)
// }
//
// func init() {
// 	rand.Seed(config.GinkgoConfig.RandomSeed)
// }
