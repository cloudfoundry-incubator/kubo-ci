package upgrade_tests_test

import (
	"testing"

	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	kubectl    *test_helpers.KubectlRunner
	iaas      string
)

func TestUpgradeTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpgradeTests Suite")
}

var _ = BeforeSuite(func() {
	kubectl = test_helpers.NewKubectlRunner()
	kubectl.Setup()
	iaas = test_helpers.GetIaas()
})

var _ = AfterSuite(func() {
	kubectl.StartKubectlCommand("delete", "--all", "psp")
	kubectl.Teardown()
})
