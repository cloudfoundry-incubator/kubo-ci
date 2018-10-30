package k8s_lbs_test

import (
	"fmt"

	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Internal load balancers", func() {
	It("exposes routes via LBs", func() {
		if iaas != "gce" {
			Skip("Test only valid for GCE")
		}
		deployNginx := runner.RunKubectlCommand("create", "-f", internalNginxLBSpec)
		Eventually(deployNginx, "60s").Should(gexec.Exit(0))
		rolloutWatch := runner.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, "120s").Should(gexec.Exit(0))

		loadbalancerAddress := ""
		Eventually(func() string {
			loadbalancerAddress = runner.GetLBAddress("nginx", iaas)
			return loadbalancerAddress
		}, "240s", "5s").Should(Not(Equal("")))

		appUrl := fmt.Sprintf("http://%s", loadbalancerAddress)

		session := runner.RunKubectlCommandInNamespace("default", "run", "test-master-cert-via-curl-"+test_helpers.GenerateRandomUUID(), "--image=appropriate/curl", "--restart=Never", "-ti", "--rm", "--", "curl", appUrl)
		Eventually(session, "5m").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		session := runner.RunKubectlCommand("delete", "-f", internalNginxLBSpec)
		session.Wait("60s")
	})
})
