package external_traffic_policy_test

import (
	"testing"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestK8sLb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ExternalTrafficPolicy Suite")
}

var (
	kubectl                *test_helpers.KubectlRunner
	echoserverLBSpec       = test_helpers.PathFromRoot("specs/echoserver-lb.yml")
	echoserverNodePortSpec = test_helpers.PathFromRoot("specs/echoserver-nodeport.yml")

	iaas string
)

var _ = BeforeSuite(func() {
	kubectl = test_helpers.NewKubectlRunner()
	kubectl.Setup()

	var err error
	iaas, err = test_helpers.IaaS()
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if kubectl != nil {
		kubectl.Teardown()
	}
})
