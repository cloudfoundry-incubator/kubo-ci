package generic_test

import (
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGeneric(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scaling Suite")
}

var runner *test_helpers.KubectlRunner

var _ = BeforeSuite(func() {
	runner = test_helpers.NewKubectlRunner()
	runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
})

var _ = AfterSuite(func() {
	if runner != nil {
		runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	}
})
