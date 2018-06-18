package k8s_lbs_test

import (
	"testing"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestK8sLb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8sLb Suite")
}

var (
	runner           *test_helpers.KubectlRunner
	nginxLBSpec      = test_helpers.PathFromRoot("specs/nginx-lb.yml")
	echoserverLBSpec = test_helpers.PathFromRoot("specs/echoserver-lb.yml")
	iaas             string
)

var _ = BeforeSuite(func() {
	runner = test_helpers.NewKubectlRunner()
	runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")

	var err error
	iaas, err = test_helpers.IaaS()
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if runner != nil {
		runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	}
})
