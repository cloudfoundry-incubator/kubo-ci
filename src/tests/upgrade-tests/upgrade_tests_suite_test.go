package upgrade_tests_test

import (
	"testing"

	"tests/config"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	k8sRunner  *test_helpers.KubectlRunner
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

	k8sRunner = test_helpers.NewKubectlRunner()
	k8sRunner.Setup()
})

var _ = AfterSuite(func() {
	k8sRunner.RunKubectlCommand("delete", "--all", "psp")
	k8sRunner.Teardown()
})
