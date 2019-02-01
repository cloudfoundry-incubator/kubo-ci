package cloud_provider_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"tests/test_helpers"
)

var kubectl *test_helpers.KubectlRunner

func TestCloudProviderTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cloud Provider Suite")
}

var _ = BeforeSuite(func() {
	kubectl = test_helpers.NewKubectlRunner()
	kubectl.Setup()
})

var _ = AfterSuite(func() {
	kubectl.Teardown()
})
