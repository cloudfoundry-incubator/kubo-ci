package gcp_lb_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"integration-tests/test_helpers"
)

func TestGcpLb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GcpLb Suite")
}

var (
	runner        *test_helpers.KubectlRunner
	nginxSpec      = test_helpers.PathFromRoot("specs/nginx.yml")
	workerAddress string
)

var _ = BeforeSuite(func() {
	workerAddress = os.Getenv("WORKER_IP_ADDRESS")

	if workerAddress == "" {
		Fail("WORKER_IP_ADDRESS is not set")
	}
	runner = test_helpers.NewKubectlRunner()
	runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
})

var _ = AfterSuite(func() {
	if runner != nil {
		runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	}
})
