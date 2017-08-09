package workload_test_k8s

import (
	"integration-tests/test_helpers"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestK8sLb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8sLb Suite")
}

var (
	runner        *test_helpers.KubectlRunner
	nginxLBSpec   = test_helpers.PathFromRoot("specs/nginx-lb.yml")
)

var _ = BeforeSuite(func() {
	runner = test_helpers.NewKubectlRunner()
	runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
})

var _ = AfterSuite(func() {
	if runner != nil {
		runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	}
})
