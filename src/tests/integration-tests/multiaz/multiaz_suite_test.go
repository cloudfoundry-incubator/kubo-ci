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
	kubectl   *test_helpers.KubectlRunner
	nginxSpec = test_helpers.PathFromRoot("specs/nginx-daemonset.yml")
)

var _ = BeforeSuite(func() {
	kubectl = test_helpers.NewKubectlRunner()
	kubectl.Setup()
})

var _ = AfterSuite(func() {
	if kubectl != nil {
		kubectl.Teardown()
	}
})
