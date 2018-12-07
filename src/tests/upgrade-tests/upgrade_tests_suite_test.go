package upgrade_tests_test

import (
	"testing"

	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var (
	kubectl *test_helpers.KubectlRunner
	iaas    string
)

func TestUpgradeTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpgradeTests Suite")
}

var _ = BeforeSuite(func() {
	test_helpers.CheckRequiredEnvs([]string{
		"BOSH_DEPLOY_COMMAND",
		"BOSH_DEPLOYMENT",
		"BOSH_ENVIRONMENT",
		"BOSH_CLIENT",
		"BOSH_CLIENT_SECRET",
		"BOSH_CA_CERT",
		"ENABLE_MULTI_AZ_TESTS",
		"IAAS",
	})

	kubectl = test_helpers.NewKubectlRunner()
	kubectl.Setup()
	iaas = os.Getenv("IAAS")
})

var _ = AfterSuite(func() {
	kubectl.StartKubectlCommand("delete", "--all", "psp")
	kubectl.Teardown()
})
