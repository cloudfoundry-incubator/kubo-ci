package upgrade_tests_test

import (
	"testing"

	"tests/config"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	kubectl    *test_helpers.KubectlRunner
	testconfig *config.Config
)

func TestUpgradeTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpgradeTests Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())

	kubectl = test_helpers.NewKubectlRunner()
	kubectl.Setup()
})

var _ = AfterSuite(func() {
	kubectl.RunKubectlCommand("delete", "--all", "psp")
	kubectl.Teardown()
})
