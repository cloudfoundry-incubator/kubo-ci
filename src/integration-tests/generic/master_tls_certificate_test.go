package generic_test

import (
	"integration-tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("MasterTlsCertificate", func() {

	var (
		runner *test_helpers.KubectlRunner
	)

	BeforeEach(func() {
		runner = test_helpers.NewKubectlRunner()
	})

	FIt("should be valid for kube-dns names for master", func() {
		session := runner.RunKubectlCommand("run", "test-curl", "--image=\"governmentpaas/curl-ssl\"", "--restart=Never", "-ti", "--rm", "--", "curl", "https://kubernetes", "--ca-cert", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
		Eventually(session).Should(gexec.Exit(0))
	})
})
