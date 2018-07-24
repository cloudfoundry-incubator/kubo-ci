package generic_test

import (
	"crypto/tls"
	"fmt"
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
		kubectl = NewKubectlRunner()
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

		podName := GenerateRandomUUID()
		session := kubectl.RunKubectlCommand("run", podName, "--image", "pcfkubo/alpine:stable", "--restart=Never", "--image-pull-policy=Always", "-ti", "--rm", "--", "kubectl", "get", "services")
		session.Wait(120)
		Expect(session).To(gexec.Exit(0))
	})

	It("Should be able to run kubectl top pod successfully", func() {
		Eventually(func() int {
			return kubectl.RunKubectlCommand("top", "pods", "--heapster-scheme=https", "-n", "kube-system").Wait(30 * time.Second).ExitCode()
		}, "300s", "10s").Should(Equal(0))
	})

	It("Should be able to run kubectl top nodes successfully", func() {
		Eventually(func() int {
			return kubectl.RunKubectlCommand("top", "nodes", "--heapster-scheme=https").Wait(30 * time.Second).ExitCode()
		}, "300s", "10s").Should(Equal(0))
	})

	Context("Dashboard", func() {
		It("Should provide access to the dashboard via kubectl proxy", func() {
			session := kubectl.RunKubectlCommand("proxy")
			Eventually(session).Should(gbytes.Say("Starting to serve on"))

			timeout := time.Duration(5 * time.Second)
			httpClient := http.Client{
				Timeout: timeout,
			}

			// For more details, see: https://github.com/kubernetes/dashboard/wiki/Accessing-Dashboard---1.7.X-and-above#kubectl-proxy
			appUrl := "http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/"

			Eventually(func() int {
				result, err := httpClient.Get(appUrl)
				if err != nil {
					return -1
				}
				return result.StatusCode
			}, "120s", "5s").Should(Equal(200))

			session.Terminate()
		})

		It("Should provide access to the dashboard via a node port", func() {

			timeout := time.Duration(5 * time.Second)
			transport := &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
			httpClient := http.Client{
				Timeout:   timeout,
				Transport: transport,
			}

			appAddress := kubectl.GetAppAddressInNamespace("svc/kubernetes-dashboard", "kube-system")
			appUrl := fmt.Sprintf("https://%s", appAddress)

			Eventually(func() int {
				result, err := httpClient.Get(appUrl)
				if err != nil {
					return -1
				}
				return result.StatusCode
			}, "120s", "5s").Should(Equal(200))
		})
	})

	Context("When unauthorized service account", func() {
		var serviceAccount string

		BeforeEach(func() {
			serviceAccount = PathFromRoot("specs/build-robot-service-account.yml")
			kubectl.RunKubectlCommandWithTimeout("create", "-f", serviceAccount)
		})

		AfterEach(func() {
			kubectl.RunKubectlCommand("delete", "-f", serviceAccount)
		})

		It("Should not be allowed to perform attach,exec,logs actions", func() {
			session := kubectl.RunKubectlCommand("--as=system:serviceaccounts:build-robot", "auth", "can-i", "attach", "pod")
			Eventually(session, "15s").Should(gbytes.Say("no"))
			session = kubectl.RunKubectlCommand("--as=system:serviceaccounts:build-robot", "auth", "can-i", "logs", "pod")
			Eventually(session, "15s").Should(gbytes.Say("no"))
			session = kubectl.RunKubectlCommand("--as=system:serviceaccounts:build-robot", "auth", "can-i", "exec", "pod")
			Eventually(session, "15s").Should(gbytes.Say("no"))
		})
	})
})
