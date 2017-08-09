package pod_logs

import (
	"integration-tests/test_helpers"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPodLogs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PodLogs Suite")
}

var (
	runner        *test_helpers.KubectlRunner
	nginxSpec     = test_helpers.PathFromRoot("specs/nginx.yml")
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
