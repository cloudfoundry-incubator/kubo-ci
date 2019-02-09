package k8s_lbs_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Internal load balancers", func() {
	It("exposes routes via LBs", func() {
		switch iaas {
		case "vsphere", "openstack":
			Skip("Test only valid for GCE, Azure and AWS")
		}
		deployNginx := kubectl.StartKubectlCommand("create", "-f", internalNginxLBSpec)
		Eventually(deployNginx, kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		rolloutWatch := kubectl.StartKubectlCommand("rollout", "status", "deployment/nginx", "-w")
		Eventually(rolloutWatch, kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))

		loadbalancerAddress := ""
		Eventually(func() string {
			loadbalancerAddress = kubectl.GetLBAddress("nginx", iaas)
			return loadbalancerAddress
		}, 15*kubectl.TimeoutInSeconds, "15s").Should(Not(Equal("")))

		appUrl := fmt.Sprintf("http://%s", loadbalancerAddress)

		Eventually(func() int {
			session := kubectl.StartKubectlCommand("run",
				"test-master-cert-via-curl",
				"--generator=run-pod/v1",
				"--image=tutum/curl",
				"--restart=Never",
				"-ti",
				"--rm",
				"--",
				"curl",
				appUrl)
			session.Wait(kubectl.TimeoutInSeconds)
			return session.ExitCode()
		}, 5*kubectl.TimeoutInSeconds, "5s").Should(Equal(0))
	})

	AfterEach(func() {
		kubectl.StartKubectlCommand("delete", "-f", internalNginxLBSpec).Wait(kubectl.TimeoutInSeconds)
	})
})
