package multiaz

import (
	"testing"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestK8sMultiAZ(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8s Multi-AZ Suite")
}

var (
	runner    *test_helpers.KubectlRunner
	nginxSpec = test_helpers.PathFromRoot("specs/nginx-daemonset.yml")
)

var _ = BeforeSuite(func() {
	runner = test_helpers.NewKubectlRunner()
	runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
})

var _ = AfterSuite(func() {
	if runner != nil && !CurrentGinkgoTestDescription().Failed {
		runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	}
})
