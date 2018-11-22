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
	runner                 *test_helpers.KubectlRunner
	echoserverLBSpec       = test_helpers.PathFromRoot("specs/echoserver-lb.yml")
	echoserverNodePortSpec = test_helpers.PathFromRoot("specs/echoserver-nodeport.yml")

	iaas string
)

var _ = BeforeSuite(func() {
	runner = test_helpers.NewKubectlRunner()
	runner.Setup()

	var err error
	iaas, err = test_helpers.IaaS()
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if runner != nil {
		runner.Teardown()
	}
})
