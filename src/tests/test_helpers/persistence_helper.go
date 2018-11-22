package test_helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func UndeployGuestBook(kubectl *KubectlRunner) {
	guestBookSpec := PathFromRoot("specs/pv-guestbook.yml")
	Eventually(kubectl.RunKubectlCommand("delete", "-f", guestBookSpec), kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
}

func DeployGuestBook(kubectl *KubectlRunner) {
	guestBookSpec := PathFromRoot("specs/pv-guestbook.yml")
	Eventually(kubectl.RunKubectlCommand("apply", "-f", guestBookSpec), 5*kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
	WaitForPodsToRun(kubectl, 5*kubectl.TimeoutInSeconds)
}

func PostToGuestBook(address string, testValue string) error {
	url := fmt.Sprintf("http://%s/guestbook.php?cmd=set&key=messages&value=%s", address, testValue)

	httpClient := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	_, err := httpClient.Get(url)

	return err
}

func GetValueFromGuestBook(address string) string {

	httpClient := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	url := fmt.Sprintf("http://%s/guestbook.php?cmd=get&key=messages", address)
	response, err := httpClient.Get(url)
	if err != nil {
		return fmt.Sprintf("error occurred : %s", err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).ToNot(HaveOccurred())
	return string(bodyBytes)

}
