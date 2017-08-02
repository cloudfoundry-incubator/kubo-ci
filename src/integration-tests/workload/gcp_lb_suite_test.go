package workload_test

import (
	"integration-tests/test_helpers"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGcpLb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GcpLb Suite")
}

var (
	runner        *test_helpers.KubectlRunner
	nginxSpec     = test_helpers.PathFromRoot("specs/nginx.yml")
	nginxLBSpec   = test_helpers.PathFromRoot("specs/nginx-lb.yml")
	workerAddress string
	nodePort      string
	iaas          string
)

var _ = BeforeSuite(func() {
	workerAddress = os.Getenv("WORKLOAD_ADDRESS")
	nodePort = os.Getenv("WORKLOAD_PORT")
	iaas = os.Getenv("IAAS")
	runner = test_helpers.NewKubectlRunner()
	runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
})

var _ = AfterSuite(func() {
	if runner != nil {
		runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	}
})
