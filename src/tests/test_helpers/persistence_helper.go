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
	timeout := time.Duration(float64(2 * time.Minute))
	Eventually(kubectl.RunKubectlCommand("delete", "-f", guestBookSpec), timeout).Should(gexec.Exit(0))
}

func DeployGuestBook(kubectl *KubectlRunner) {
	guestBookSpec := PathFromRoot("specs/pv-guestbook.yml")
	timeout := "5m"
	Eventually(kubectl.RunKubectlCommand("apply", "-f", guestBookSpec), timeout).Should(gexec.Exit(0))
	WaitForPodsToRun(kubectl, timeout)
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
