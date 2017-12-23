package generic_test

import (
	. "tests/test_helpers"

	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Kubectl", func() {
	var (
		kubectl *KubectlRunner
	)

	BeforeEach(func() {
		kubectl = NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
		kubectl.RunKubectlCommand(
			"create", "namespace", kubectl.Namespace()).Wait("60s")
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand(
			"delete", "namespace", kubectl.Namespace()).Wait("60s")
	})

	It("Should be able to run kubectl commands within pod", func() {

		roleBindingName := kubectl.Namespace() + "-admin"
		s := kubectl.RunKubectlCommand("create", "rolebinding", roleBindingName, "--clusterrole=admin", "--user=system:serviceaccount:"+kubectl.Namespace()+":default")
		Eventually(s, "15s").Should(gexec.Exit(0))

		podName := GenerateRandomName()
		session := kubectl.RunKubectlCommand("run", podName, "--image", "pcfkubo/alpine:stable", "--restart=Never", "--image-pull-policy=Always", "-ti", "--rm", "--", "kubectl", "get", "services")
		session.Wait(120)
		Expect(session).To(gexec.Exit(0))
	})

	It("Should provide access to the dashboard", func() {
		session := kubectl.RunKubectlCommand("proxy")
		Eventually(session).Should(gbytes.Say("Starting to serve on"))

		timeout := time.Duration(5 * time.Second)
		httpClient := http.Client{
			Timeout: timeout,
		}

		appUrl := "http://127.0.0.1:8001/ui/"

		Eventually(func() int {
			result, err := httpClient.Get(appUrl)
			if err != nil {
				return -1
			}
			return result.StatusCode
		}, "120s", "5s").Should(Equal(200))

		session.Terminate()
	})

})
