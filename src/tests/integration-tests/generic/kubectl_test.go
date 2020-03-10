package generic_test

import (
	. "tests/test_helpers"

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
		kubectl.Setup()
	})

	AfterEach(func() {
		kubectl.Teardown()
	})

	It("Should be able to run kubectl commands within pod", func() {
		roleBindingName := kubectl.Namespace() + "-admin"
		s := kubectl.StartKubectlCommand("create", "rolebinding", roleBindingName, "--clusterrole=admin", "--user=system:serviceaccount:"+kubectl.Namespace()+":default")
		Eventually(s, "15s").Should(gexec.Exit(0))

		podName := GenerateRandomUUID()
		session := kubectl.StartKubectlCommand("run", podName, "--image", "pcfkubo/alpine:stable", "--restart=Never", "--image-pull-policy=Always", "-ti", "--rm", "--", "kubectl", "get", "services")
		session.Wait(120)
		Expect(session).To(gexec.Exit(0))
	})

	It("Should be able to run kubectl top pod successfully", func() {
		Eventually(func() int {
			return kubectl.StartKubectlCommand("top", "pods", "-n", "kube-system").Wait(30 * time.Second).ExitCode()
		}, "300s", "10s").Should(Equal(0))
	})

	It("Should be able to run kubectl top nodes successfully", func() {
		Eventually(func() int {
			return kubectl.StartKubectlCommand("top", "nodes").Wait(30 * time.Second).ExitCode()
		}, "300s", "10s").Should(Equal(0))
	})

	Context("When unauthorized service account", func() {
		var serviceAccount string

		BeforeEach(func() {
			serviceAccount = PathFromRoot("specs/build-robot-service-account.yml")
			kubectl.RunKubectlCommandWithTimeout("create", "-f", serviceAccount)
		})

		AfterEach(func() {
			kubectl.StartKubectlCommand("delete", "-f", serviceAccount).Wait(kubectl.TimeoutInSeconds)
		})

		It("Should not be allowed to perform attach,exec,logs actions", func() {
			session := kubectl.StartKubectlCommand("--as=system:serviceaccounts:build-robot", "auth", "can-i", "attach", "pod")
			Eventually(session, "15s").Should(gbytes.Say("no"))
			session = kubectl.StartKubectlCommand("--as=system:serviceaccounts:build-robot", "auth", "can-i", "logs", "pod")
			Eventually(session, "15s").Should(gbytes.Say("no"))
			session = kubectl.StartKubectlCommand("--as=system:serviceaccounts:build-robot", "auth", "can-i", "exec", "pod")
			Eventually(session, "15s").Should(gbytes.Say("no"))
		})
	})
})
