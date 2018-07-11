package generic_test

import (
	"fmt"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("MasterTlsCertificate", func() {

	var (
		kubectl *KubectlRunner
	)

	BeforeEach(func() {
		kubectl = NewKubectlRunner()
	})

	DescribeTable("hostnames", func(hostname string) {
		url := fmt.Sprintf("https://%s", hostname)

		Eventually(func() *gexec.Session {
			session := kubectl.RunKubectlCommandInNamespace("default", "run", "test-master-cert-via-curl-"+GenerateRandomUUID(),
				"--image=tutum/curl", "--restart=Never", "-ti", "--rm", "--",
				"curl", url, "--cacert", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
			session.Wait("5s")
			return session
		}, "1m", "10s").Should(gexec.Exit(0))
	},
		Entry("kubernetes", "kubernetes"),
		Entry("kubernetes.default", "kubernetes.default"),
		Entry("kubernetes.default.svc", "kubernetes.default.svc"),
		Entry("kubernetes.default.svc.cluster.local", "kubernetes.default.svc.cluster.local"),
	)
})
