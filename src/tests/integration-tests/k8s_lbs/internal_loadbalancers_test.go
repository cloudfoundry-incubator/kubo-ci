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
		deployNginx := kubectl.RunKubectlCommand("create", "-f", internalNginxLBSpec)
		Eventually(deployNginx, kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		rolloutWatch := kubectl.RunKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))

		loadbalancerAddress := ""
		Eventually(func() string {
			loadbalancerAddress = kubectl.GetLBAddress("nginx", iaas)
			return loadbalancerAddress
		}, "240s", "5s").Should(Not(Equal("")))

		appUrl := fmt.Sprintf("http://%s", loadbalancerAddress)

		session := kubectl.RunKubectlCommandInNamespace("default", "run", "test-master-cert-via-curl-"+test_helpers.GenerateRandomUUID(), "--image=tutum/curl", "--restart=Never", "-ti", "--rm", "--", "curl", appUrl)
		Eventually(session, "5m").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "-f", internalNginxLBSpec).Wait(kubectl.TimeoutInSeconds)
	})
})
