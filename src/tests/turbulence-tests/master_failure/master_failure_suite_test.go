package master_failure_test

import (
	"tests/config"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	iaas       string
	testconfig *config.Config
)

func TestMasterFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Master failure suite")
}

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
	director := NewDirector(testconfig.Bosh)
	director.EnableResurrection(true)
})

func MasterFailureDescribe(description string, callback func()) bool {
	return Describe("[master_failure]", func() {
		BeforeEach(func() {
			if !testconfig.TurbulenceTests.IncludeMasterFailure {
				Skip(`Skipping this test suite because Config.TurbulenceTests.IncludeMasterFailure is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
