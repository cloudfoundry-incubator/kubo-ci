package test_helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func UndeployGuestBook(kubectl *KubectlRunner, timeoutScale float64) {
	guestBookSpec := PathFromRoot("specs/pv-guestbook.yml")
	timeout := time.Duration(float64(2*time.Minute) * timeoutScale)
	Eventually(kubectl.RunKubectlCommand("delete", "-f", guestBookSpec), timeout).Should(gexec.Exit(0))
}

func DeployGuestBook(kubectl *KubectlRunner, timeoutScale float64) {
	guestBookSpec := PathFromRoot("specs/pv-guestbook.yml")
	timeout := time.Duration(float64(2*time.Minute) * timeoutScale)
	Eventually(kubectl.RunKubectlCommand("apply", "-f", guestBookSpec), timeout).Should(gexec.Exit(0))
	Eventually(func() *gexec.Session {
		session := kubectl.RunKubectlCommand("rollout", "status", "deployment/frontend", "--watch=false")
		session.Wait(2 * time.Minute)
		return session
	}, timeout, 10*time.Second).Should(gbytes.Say("successfully rolled out"))
	Eventually(func() *gexec.Session {
		session := kubectl.RunKubectlCommand("rollout", "status", "deployment/redis-master", "--watch=false")
		session.Wait(2 * time.Minute)
		return session
	}, timeout, 10*time.Second).Should(gbytes.Say("successfully rolled out"))
}

func PostToGuestBook(address string, testValue string) {

	url := fmt.Sprintf("http://%s/guestbook.php?cmd=set&key=messages&value=%s", address, testValue)
	_, err := http.Get(url)
	Expect(err).ToNot(HaveOccurred())

}

func GetValueFromGuestBook(address string) string {

	httpClient := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	url := fmt.Sprintf("http://%s/guestbook.php?cmd=get&key=messages", address)
	response, err := httpClient.Get(url)
	if err != nil {
		return fmt.Sprintf("error occured : %s", err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).ToNot(HaveOccurred())
	return string(bodyBytes)

}
