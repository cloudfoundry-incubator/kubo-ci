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

var kubectl *test_helpers.KubectlRunner

var _ = BeforeSuite(func() {
	test_helpers.CheckRequiredEnvs([]string{"HPA_TIMEOUT"})
	kubectl = test_helpers.NewKubectlRunner()
	kubectl.Setup()
})

var _ = AfterSuite(func() {
	if kubectl != nil {
		kubectl.Teardown()
	}
})
